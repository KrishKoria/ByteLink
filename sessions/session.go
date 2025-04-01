package sessions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/KrishKoria/ByteLink/internal/database"
	"github.com/KrishKoria/ByteLink/models"
	"github.com/google/uuid"
	"time"
)

var (
	authService *models.AuthService
	ctx         = context.Background()
)

func CreateSession(userID string) (*models.Session, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(authService.SessionTimeout)

	err := authService.DB.CreateSession(ctx, database.CreateSessionParams{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &models.Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}, nil
}

func ValidateSession(token string) (*models.User, error) {
	session, err := authService.DB.GetSession(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid session")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	if session.ExpiresAt.Before(time.Now()) {
		err := authService.DB.DeleteSession(ctx, token)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("session expired")
	}

	user, err := authService.DB.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
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

func EndSession(token string) error {
	return authService.DB.DeleteSession(ctx, token)

}
