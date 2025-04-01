package handler

import (
	"context"
	"fmt"
	"github.com/KrishKoria/ByteLink/miscellaneous"
	"github.com/KrishKoria/ByteLink/shortener"
	"github.com/KrishKoria/ByteLink/store"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := store.IncrementClickCount(ctx, shortUrl)
		if err != nil {
			fmt.Printf("Failed to increment click count: %v\n", err)
		}
	}()
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

func DeleteUserURL(c *gin.Context) {
	shortURL := c.Query("short_url")
	userID := c.Query("user_id")

	if shortURL == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "short_url and user_id parameters are required",
		})
		return
	}

	err := store.DeleteMapping(shortURL, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "URL mapping deleted successfully",
	})
}

func GetCleanupStatus(c *gin.Context) {
	status := miscellaneous.CleanupStatus{
		LastRunTime:        time.Now(),
		TotalURLsRemoved:   0,
		IsRunning:          true,
		RunIntervalMinutes: 60,
	}

	c.JSON(http.StatusOK, status)
}

func GetURLStatsHandler(c *gin.Context) {
	shortURL := c.Query("short_url")
	userID := c.Query("user_id")

	if shortURL == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters",
		})
		return
	}

	stats, err := store.GetURLStats(c, shortURL, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stats not found"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
