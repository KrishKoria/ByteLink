package models

import (
	"github.com/KrishKoria/ByteLink/internal/database"
	"golang.org/x/oauth2"
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	Password      string    `json:"password,omitempty"` // Omit in responses
	Provider      string    `json:"provider"`           // "local", "google", "github"
	ProviderID    string    `json:"-"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type AuthService struct {
	DB             *database.Queries
	GoogleConfig   *oauth2.Config
	GithubConfig   *oauth2.Config
	SessionTimeout time.Duration
}
