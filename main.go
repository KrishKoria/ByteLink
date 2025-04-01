package main

import (
	"fmt"
	"github.com/KrishKoria/ByteLink/authentication"
	"github.com/KrishKoria/ByteLink/handler"
	"github.com/KrishKoria/ByteLink/internal/database"
	"github.com/KrishKoria/ByteLink/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	store.InitializeStoreService()
	authentication.InitializeAuthService(database.New(store.GetDBConn()))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to ByteLink",
		})
	})

	// Authentication routes
	authGroup := r.Group("/auth")
	{
		// Traditional email/password authentication
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/logout", handler.Logout)

		// OAuth routes
		authGroup.GET("/google", handler.GoogleLogin)
		authGroup.GET("/google/callback", handler.GoogleCallback)
		authGroup.GET("/github", handler.GithubLogin)
		authGroup.GET("/github/callback", handler.GithubCallback)
	}

	// Public routes
	r.POST("/create", func(c *gin.Context) {
		handler.CreateShortURL(c)
	})

	r.GET("/:shortURL", func(c *gin.Context) {
		handler.HandleRedirect(c)
	})

	// Protected routes
	apiGroup := r.Group("/api")
	apiGroup.Use(authentication.AuthMiddleware())
	{
		apiGroup.GET("/url", handler.GetUserURL)
		apiGroup.GET("/urls", handler.GetUserURLs)
		apiGroup.DELETE("/url", handler.DeleteUserURL)
		apiGroup.GET("/url/stats", handler.GetURLStatsHandler)

		// Admin routes
		apiGroup.GET("/admin/cleanup-status", handler.GetCleanupStatus)
	}

	err := r.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
