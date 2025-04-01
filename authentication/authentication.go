package authentication

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/KrishKoria/ByteLink/internal/database"
	"github.com/KrishKoria/ByteLink/models"
	"github.com/KrishKoria/ByteLink/sessions"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"io"
	"os"
	"time"
)

var (
	AuthService *models.AuthService
	ctx         = context.Background()
)

func InitializeAuthService(db *database.Queries) {
	AuthService = &models.AuthService{
		DB: db,
		GoogleConfig: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("BASE_URL") + "/auth/google/callback",
			Scopes:       []string{"profile", "email"},
			Endpoint:     google.Endpoint,
		},
		GithubConfig: &oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("BASE_URL") + "/auth/github/callback",
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
		SessionTimeout: 24 * time.Hour,
	}
}

func RegisterUser(email, name, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	passwordHash := string(hashedPassword)
	userID := uuid.New().String()
	now := time.Now()

	user, err := AuthService.DB.CreateUser(ctx, database.CreateUserParams{
		ID:            userID,
		Email:         email,
		Name:          name,
		Password:      sql.NullString{String: passwordHash, Valid: true},
		Provider:      "local",
		EmailVerified: sql.NullBool{Bool: true, Valid: true},
		CreatedAt:     now,
		UpdatedAt:     now,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &models.User{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		Provider:      user.Provider,
		EmailVerified: user.EmailVerified.Bool,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}

func LoginUser(email, password string) (*models.Session, error) {
	user, err := AuthService.DB.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.Provider != "local" {
		return nil, fmt.Errorf("this account uses %s authentication", user.Provider)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return sessions.CreateSession(user.ID)
}

func GetGoogleAuthURL(state string) string {
	return AuthService.GoogleConfig.AuthCodeURL(state)
}

func GetGithubAuthURL(state string) string {
	return AuthService.GithubConfig.AuthCodeURL(state)
}

func HandleGoogleCallback(code string) (*models.Session, error) {
	token, err := AuthService.GoogleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	client := AuthService.GoogleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// TODO : Parse user info and handle user creation/login
	// (simplified - you'd need to parse the JSON response)
	// This is pseudocode - you need to implement the actual JSON parsing
	var googleUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		ID    string `json:"id"`
	}
	// Parse JSON into googleUser

	return handleOAuthUser(googleUser.Email, googleUser.Name, "google", googleUser.ID)
}
func HandleGithubCallback(code string) (*models.Session, error) {
	token, err := AuthService.GithubConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	client := AuthService.GithubConfig.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// TODO: Parse user info and handle user creation/login
	// (simplified - you'd need to parse the JSON response)
	// This is pseudocode - you need to implement the actual JSON parsing
	var githubUser struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	// Parse JSON into githubUser

	// For GitHub, you might need to make another request to get email if not present

	return handleOAuthUser(githubUser.Email, githubUser.Name, "github", fmt.Sprintf("%d", githubUser.ID))
}

func handleOAuthUser(email, name, provider, providerID string) (*models.Session, error) {
	user, err := AuthService.DB.GetUserByEmail(ctx, email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		userID := uuid.New().String()
		now := time.Now()

		user, err = AuthService.DB.CreateOAuthUser(ctx, database.CreateOAuthUserParams{
			ID:            userID,
			Email:         email,
			Name:          name,
			Provider:      provider,
			ProviderID:    sql.NullString{String: providerID, Valid: true},
			EmailVerified: sql.NullBool{Bool: true, Valid: true},
			CreatedAt:     now,
			UpdatedAt:     now,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else if user.Provider != provider {
		return nil, fmt.Errorf("this email is already used with %s authentication", user.Provider)
	}

	return sessions.CreateSession(user.ID)
}
