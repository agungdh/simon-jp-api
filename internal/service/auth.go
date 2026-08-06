package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"simon-jp-api/internal/models"
	"simon-jp-api/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrSessionNotFound    = errors.New("session not found")
)

type AuthService struct {
	users   *repository.UserRepository
	session *SessionStore
}

func NewAuthService(users *repository.UserRepository, session *SessionStore) *AuthService {
	return &AuthService{users: users, session: session}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, time.Duration, error) {
	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return "", 0, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", 0, ErrInvalidCredentials
	}

	token, err := generateToken()
	if err != nil {
		return "", 0, fmt.Errorf("generate token: %w", err)
	}

	if err := s.session.Create(ctx, token, user.ID); err != nil {
		return "", 0, fmt.Errorf("create session: %w", err)
	}

	return token, s.session.TTL(), nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	if _, err := s.session.Get(ctx, token); err != nil {
		return ErrInvalidToken
	}
	return s.session.Delete(ctx, token)
}

func (s *AuthService) Me(ctx context.Context, token string) (*models.User, error) {
	userID, err := s.session.Get(ctx, token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return user, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
