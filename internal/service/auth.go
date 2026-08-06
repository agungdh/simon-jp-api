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
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrSessionNotFound    = errors.New("session not found")
)

type UserStore interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByID(ctx context.Context, id int64) (*models.User, error)
}

type PegawaiStore interface {
	FindByUserID(ctx context.Context, userID int64) (*models.Pegawai, error)
}

type LoginResult struct {
	Token     string
	ExpiresIn time.Duration
	User      *models.User
	Pegawai   *models.Pegawai
}

type AuthService struct {
	users    UserStore
	pegawai  PegawaiStore
	session  SessionStore
	throttle Throttler
}

func NewAuthService(users UserStore, pegawai PegawaiStore, session SessionStore, throttle Throttler) *AuthService {
	return &AuthService{users: users, pegawai: pegawai, session: session, throttle: throttle}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	key := "login:" + username
	if err := s.throttle.Check(ctx, key); err != nil {
		return nil, err
	}

	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		dummyHash, _ := bcrypt.GenerateFromPassword([]byte("dummy"), bcrypt.MinCost)
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
		_ = s.throttle.RecordFailure(ctx, key)
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		_ = s.throttle.RecordFailure(ctx, key)
		return nil, ErrInvalidCredentials
	}

	_ = s.throttle.Reset(ctx, key)

	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	if err := s.session.Create(ctx, token, user.ID); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	pegawai, _ := s.pegawai.FindByUserID(ctx, user.ID)

	return &LoginResult{
		Token:     token,
		ExpiresIn: s.session.TTL(),
		User:      user,
		Pegawai:   pegawai,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	if _, err := s.session.Get(ctx, token); err != nil {
		return ErrInvalidToken
	}
	return s.session.Delete(ctx, token)
}

func (s *AuthService) Authenticate(ctx context.Context, token string) (*models.User, error) {
	userID, err := s.session.Get(ctx, token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return user, nil
}

func (s *AuthService) Pegawai(ctx context.Context, userID int64) (*models.Pegawai, error) {
	pegawai, err := s.pegawai.FindByUserID(ctx, userID)
	if err != nil {
		return nil, nil
	}
	return pegawai, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
