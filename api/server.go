package api

import (
	"context"
	"fmt"
	"github.com/Sem4kok/restful-api/db"
	"github.com/Sem4kok/restful-api/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

var albums []util.Album

func StartServer() {
	defer func() {
		_ = conn.Close(context.Background())
	}()

	albums, err = db.GetAlbumsFromDB(context.Background(), conn)
	if err != nil && albums == nil {
		log.Fatal(err)
	} else if err != nil {
		fmt.Println("error, with rows.Err(): ", err.Error())
	}

}
