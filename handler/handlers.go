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
	err := store.SaveMapping(shortUrl, request.InitialURL, request.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	host := "http://localhost:8080/"
	c.JSON(http.StatusOK, gin.H{
		"message":   "Short URL created successfully",
		"short_url": host + shortUrl,
	})
}

func HandleRedirect(c *gin.Context) {
	shortUrl := c.Param("shortURL")
	initialURL := store.GetLongUrlPublic(shortUrl)
	if initialURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.Redirect(http.StatusFound, initialURL)
}

func GetUserURL(c *gin.Context) {
	shortUrl := c.Query("short_url")
	userId := c.Query("user_id")

	if shortUrl == "" || userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters",
		})
		return
	}

	longUrl := store.GetLongUrl(shortUrl, userId)
	if longUrl == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "URL not found or not authorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_url": shortUrl,
		"long_url":  longUrl,
	})
}

func GetUserURLs(c *gin.Context) {
	userId := c.Query("user_id")
	urls := store.GetMappingsByUserID(userId)
	c.JSON(http.StatusOK, gin.H{"urls": urls})
}
