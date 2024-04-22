package api

import (
	"context"
	"github.com/Sem4kok/restful-api/db"
	"github.com/Sem4kok/restful-api/util"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	HOST = "localhost:8080"
)

type Handler struct {
	Conn *pgx.Conn
}

var albums []util.Album

func StartServer() {

	handler := &Handler{Conn: db.StartDBConnection()}
	defer func() {
		_ = handler.Conn.Close(context.Background())
	}()

	router := gin.Default()
	router.GET("/albums", handler.getAlbums)
	// post method won't update data in db
	router.POST("/albums", handler.postAlbums)
	router.DELETE("/albums/:id", handler.deleteAlbumById)

	err := router.Run(HOST)
	if err != nil {
		log.Fatal(err)
	}
}

// private method of Handler struct
// that implements Get Method
func (h *Handler) getAlbums(c *gin.Context) {

	if albums != nil {
		c.IndentedJSON(http.StatusOK, albums)
		return
	}

	var err error = nil
	albums, err = db.GetAlbumsFromDB(context.Background(), h.Conn)
	if err != nil && albums == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		log.Printf("error, with rows.Err(): %v", err.Error())
	}

	c.IndentedJSON(http.StatusOK, albums)
}

// private method of Handler struct
// that implements Post Method
func (h *Handler) postAlbums(c *gin.Context) {

	var newAlbums = make([]util.Album, 1)
	if err := c.BindJSON(&newAlbums); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get current data from db
	if albums == nil {
		var err error = nil
		albums, err = db.GetAlbumsFromDB(context.Background(), h.Conn)
		if err != nil && albums == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			log.Printf("error, with rows.Err(): %v", err.Error())
		}
	}

	// validity of sent json data is guaranteed
	// albums will automatically update
	db.AddData(&db.Config{
		Ctx:    context.Background(),
		Conn:   h.Conn,
		Albums: &albums,
	}, newAlbums)
	time.Sleep(time.Millisecond * 10)

	c.Status(http.StatusOK)
}

// deleteAlbum has O(log(n))
// Time complexity of deleting
// where n -> len(albums)
func (h *Handler) deleteAlbumById(c *gin.Context) {
	// getting id value from client
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("error. id is not integer: %v", err.Error())

	}

	// delete album from DataBase
	didExist, err := db.DeleteAlbumFromDB(&db.Config{
		Ctx:    context.Background(),
		Conn:   h.Conn,
		Albums: &albums,
	}, id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if didExist {
		deleteFromAlbums(id)
	}

	c.JSON(http.StatusOK, gin.H{"success": "album has been deleted"})
}

func deleteFromAlbums(id int) {
	var l, r int
	l = 0
	switch albums {
	case nil:
		r = 0
	default:
		r = len(albums) - 1
	}
	// delete album from slice using binary search
	for l < r {
		mid := (r-l)/2 + l
		if albums[mid].ID < id {
			l = mid + 1
		} else if albums[mid].ID > id {
			r = mid - 1
		} else {
			albums = append(albums[:mid], albums[mid+1:]...)
			break
		}
	}
}
