package main

import (
	"net/http"

	"github.com/michaellazuardi/newtodo/tree/main/api"

	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", Index)

	router.GET("/items")

	return router
}

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello worldsssss"})
}

func main() {
	r := SetRouter()
	r.Run()
}
