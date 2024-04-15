package api

import (
	"context"
	"fmt"
	"github.com/Sem4kok/restful-api/db"
	"github.com/Sem4kok/restful-api/util"
	"github.com/gin-gonic/gin"
	"log"
)

const (
	HOST = "localhost:8080"
)

var albums []util.Album

func StartServer() {

	router := gin.Default()
	router.GET("/albums", getAlbums)

	conn := db.StartDBConnection()
	defer func() {
		_ = conn.Close(context.Background())
	}()

	albums, err := db.GetAlbumsFromDB(context.Background(), conn)
	if err != nil && albums == nil {
		log.Fatal(err)
	} else if err != nil {
		fmt.Println("error, with rows.Err(): ", err.Error())
	}

	err = router.Run(HOST)
	if err != nil {
		log.Fatal(err)
	}
}
