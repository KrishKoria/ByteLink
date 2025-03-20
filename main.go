package main

import (
	"fmt"
	"github.com/KrishKoria/ByteLink/handler"
	"github.com/KrishKoria/ByteLink/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	store.InitializeStoreService()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to ByteLink",
		})
	})

	r.POST("/create", func(c *gin.Context) {
		handler.CreateShortURL(c)
	})

	r.GET("/:shortURL", func(c *gin.Context) {
		handler.HandleRedirect(c)
	})

	r.GET("/api/url", func(c *gin.Context) {
		handler.GetUserURL(c)
	})

	r.GET("/api/urls", func(c *gin.Context) {
		handler.GetUserURLs(c)
	})

	r.GET("/api/admin/cleanup-status", func(c *gin.Context) {
		c.JSON(http.StatusOK, handler.GetCleanupStatus())
	})

	r.DELETE("/api/url", func(c *gin.Context) {
		handler.DeleteUserURL(c)
	})

	r.GET("/api/url/stats", func(c *gin.Context) {
		handler.GetURLStatsHandler(c)
	})

	err := r.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
