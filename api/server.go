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
	}

	var err error = nil
	albums, err = db.GetAlbumsFromDB(context.Background(), h.Conn)
	if err != nil && albums == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal(err)
	} else if err != nil {
		log.Printf("error, with rows.Err(): %v", err.Error())
	}

	c.IndentedJSON(http.StatusOK, albums)
}
