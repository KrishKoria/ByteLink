package handler

import (
	"encoding/json"
	"github.com/KrishKoria/ByteLink/authentication"
	"github.com/KrishKoria/ByteLink/models"
	"github.com/KrishKoria/ByteLink/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := authentication.RegisterUser(req.Email, req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := authentication.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set cookie and return token
	maxAge := int(session.ExpiresAt.Sub(time.Now()).Seconds())
	c.SetCookie("session_token", session.Token, maxAge, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token":   session.Token,
		"user_id": session.UserID,
		"expires": session.ExpiresAt,
	})
}

func Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		cookie, err := c.Request.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No session found"})
			return
		}
		token = cookie.Value
	}

	if err := sessions.EndSession(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GoogleLogin initiates Google OAuth flow
func GoogleLogin(c *gin.Context) {
	// Use a cryptographically secure random string in production
	state := "random-state"
	url := authentication.GetGoogleAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the Google OAuth callback
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return
	}

	session, err := authentication.HandleGoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	maxAge := int(session.ExpiresAt.Sub(time.Now()).Seconds())
	c.SetCookie("session_token", session.Token, maxAge, "/", "", false, true)

	// Redirect to frontend or return token
	c.Redirect(http.StatusFound, "/")
}

// GithubLogin initiates GitHub OAuth flow
func GithubLogin(c *gin.Context) {
	// Use a cryptographically secure random string in production
	state := "random-state"
	url := authentication.GetGithubAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GithubCallback handles the GitHub OAuth callback
func GithubCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return
	}

	session, err := authentication.HandleGithubCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	maxAge := int(session.ExpiresAt.Sub(time.Now()).Seconds())
	c.SetCookie("session_token", session.Token, maxAge, "/", "", false, true)

	// Redirect to frontend or return token
	c.Redirect(http.StatusFound, "/")
}

// Implement the actual JSON parsing for GitHub and Google OAuth
func parseGoogleUserInfo(body io.Reader) (*models.User, error) {
	var googleUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		ID    string `json:"id"`
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &googleUser); err != nil {
		return nil, err
	}

	return &models.User{
		Email: googleUser.Email,
		Name:  googleUser.Name,
	}, nil
}

func parseGithubUserInfo(body io.Reader) (*models.User, error) {
	var githubUser struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &githubUser); err != nil {
		return nil, err
	}

	// Use login as name if name is empty
	name := githubUser.Name
	if name == "" {
		name = githubUser.Login
	}

	return &models.User{
		Email: githubUser.Email,
		Name:  name,
	}, nil
}
