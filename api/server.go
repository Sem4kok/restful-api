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
	router.GET("/albums/:id", handler.getAlbumById)
	router.POST("/albums", handler.postAlbums)
	router.DELETE("/albums/:id", handler.deleteAlbumById)
	router.PATCH("/albums/:id", handler.updateByID)

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
// postAlbums method won't update data in db
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

// getAlbumById return needed album to client in JSON format
// return error message to client if album does not exist by this id
func (h *Handler) getAlbumById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("error. id is not integer: %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}

	index := h.findElementInSlice(id, c)
	if index == -1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no such id"})
	}

	c.IndentedJSON(http.StatusOK, albums[index])
}

// deleteAlbum has O(log(n))
// Time complexity of deleting
// where n -> len(albums)
func (h *Handler) deleteAlbumById(c *gin.Context) {
	// getting id value from client
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("error. id is not integer: %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
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

// deleteFromAlbums deleting data from cache storage
// does not touch DB
func deleteFromAlbums(id int) {
	if albums == nil {
		return
	}
	l, r := 0, len(albums)-1
	// delete album from slice using binary search
	for l <= r {
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

// create albums cache if it hasn't created before
// return index in albums of searching element by ID
func (h *Handler) findElementInSlice(id int, c *gin.Context) int {
	// if albums hasn't created before
	if albums == nil {
		var err error = nil
		albums, err = db.GetAlbumsFromDB(context.Background(), h.Conn)
		if err != nil && albums == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else if err != nil {
			log.Printf("error, with rows.Err(): %v", err.Error())
		}
	}

	if len(albums) == 0 {
		return -1
	}

	l, r := 0, len(albums)-1
	for l <= r {
		mid := (r-l)/2 + l
		if albums[mid].ID < id {
			l = mid + 1
		} else if albums[mid].ID > id {
			r = mid
		} else {
			return mid
		}
	}

	return -1
}

// updateByID receive JSON album-data.
// Client can change everything but not ID
func (h *Handler) updateByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("error. id is not integer: %v", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}

	index := h.findElementInSlice(id, c)
	if index == -1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no such id"})
	}

	// update data into cache storage
	var newAlbum = &util.Album{}
	if err := c.BindJSON(&newAlbum); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newAlbum.ID = albums[index].ID
	albums[index] = *newAlbum
	log.Printf("%v", newAlbum)

	db.UpdateData(&db.AlbumConfig{
		Cfg: &db.Config{
			Ctx:  context.Background(),
			Conn: h.Conn,
		},
		NewAlbum: *newAlbum,
	})

	c.JSON(http.StatusOK, gin.H{"success": "album has been updated"})
}
