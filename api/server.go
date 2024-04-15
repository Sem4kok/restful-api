package api

import (
	"context"
	"fmt"
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
	var err error = nil
	albums, err = db.GetAlbumsFromDB(context.Background(), h.Conn)
	if err != nil && albums == nil {
		log.Fatal(err)
	} else if err != nil {
		fmt.Println("error, with rows.Err(): ", err.Error())
	}

	c.IndentedJSON(http.StatusOK, albums)
}
