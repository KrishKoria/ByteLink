package main

import (
	"fmt"
	"github.com/KrishKoria/ByteLink/handler"
	"github.com/KrishKoria/ByteLink/store"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.StaticFile("/", "./index.html")

	r.POST("/create", func(c *gin.Context) {
		handler.CreateShortURL(c)
	})

	r.GET("/:shortURL", func(c *gin.Context) {
		handler.HandleRedirect(c)
	})

	store.InitializeStoreService()

	err := r.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
