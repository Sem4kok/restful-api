package api

import (
	"context"
	"github.com/Sem4kok/restful-api/db"
	"github.com/Sem4kok/restful-api/util"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
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
	router.POST("/albums", handler.postAlbums)

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

func (h *Handler) postAlbums(c *gin.Context) {

	var newAlbums []util.Album
	if err := c.BindJSON(newAlbums); err != nil {
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

	if err := db.AddData(db.AddConfig{
		Ctx:       context.Background(),
		Conn:      h.Conn,
		Albums:    albums,
		NewAlbums: newAlbums,
	}); err != nil {
		// TODO error catching implementation
	}

	c.Status(http.StatusOK)
}
