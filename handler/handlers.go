package handler

import (
	"github.com/KrishKoria/ByteLink/shortener"
	"github.com/KrishKoria/ByteLink/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreationRequest struct {
	InitialURL string `json:"long_url" binding:"required"`
	UserID     string `json:"user_id" binding:"required"`
}

func CreateShortURL(c *gin.Context) {
	var request CreationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	shortUrl := shortener.GenerateShortURL(request.InitialURL, request.UserID)
	store.SaveMapping(shortUrl, request.InitialURL, request.UserID)
	host := "http://localhost:8080/"
	c.JSON(http.StatusOK, gin.H{
		"message":   "Short URL created successfully",
		"short_url": host + shortUrl,
	})
}

func HandleRedirect(c *gin.Context) {
	shortUrl := c.Param("shortURL")
	initialURL := store.GetLongUrl(shortUrl)
	c.Redirect(302, initialURL)
}
