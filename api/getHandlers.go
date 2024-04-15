package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
