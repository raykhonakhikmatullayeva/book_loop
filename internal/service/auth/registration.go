package auth

import (
	"book_loop/internal/model"
	"book_loop/internal/repository/auth"
	"book_loop/internal/repository/users"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenExp    = 15 * time.Minute
	refreshBuf        = 32
	passwordMinLength = 8
	loginMinLength    = 3
)

type service struct {
	userRepo  users.Repo
	authRepo  auth.Repo
	jwtSecret []byte
}

type Service interface {
	Register(ctx context.Context, login, password, role string) error
	Login(ctx context.Context, login, password string) (model.Token, error)
	Refresh(ctx context.Context, rawToken string) (model.Token, error)
	newAccessToken(userID int64, role string) (string, error)
	newRefreshToken(ctx context.Context, userID int64) (string, error)
}

func NewService(userRepo users.Repo, authRepo auth.Repo, jwtSecret []byte) Service {
	return &service{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
		authRepo: authRepo,}
}

func (s *service) Register(ctx context.Context, login, password, role string) error {
	if role == "" {
		role = "user"
	}
	if len(password) < passwordMinLength {
		return model.ErrPasswordLength
	}
	if len(login) < loginMinLength {
		return model.ErrLoginLength
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	err = s.userRepo.CreateUser(ctx, login, string(passHash), role)
	if err != nil {
		return fmt.Errorf("register user: %w", err)
	}
	return nil
}

func (s *service) Login(ctx context.Context, login, password string) (model.Token, error) {
	user, err := s.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return model.Token{}, fmt.Errorf("get user by login: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return model.Token{}, fmt.Errorf("compare password: %w", err)
	}
	accessToken, err := s.newAccessToken(user.Id, user.Role)
	if err != nil {
		return model.Token{}, fmt.Errorf("create access token: %w", err)
	}
	refreshToken, err := s.newRefreshToken(ctx, user.Id)
	if err != nil {
		return model.Token{}, fmt.Errorf("create refresh token: %w", err)
	}
	return model.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *service) newAccessToken(userID int64, role string) (string, error) {
	claims := jwt.MapClaims{
		"uid":  userID,
		"role": role,
		"exp":  time.Now().Add(accessTokenExp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}
	return tokenString, nil
}
func generateRefreshToken() (string, error) {
	buf := make([]byte, refreshBuf)
	_, err := rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func (s *service) newRefreshToken(ctx context.Context, userID int64) (string, error) {
	raw, err := generateRefreshToken()
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	hash := hashToken(raw)
	err = s.authRepo.SaveRefresh(ctx, userID, hash, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return "", fmt.Errorf("save refresh token: %w", err)
	}
	return raw, nil
}

func (s *service) Refresh(ctx context.Context, rawToken string) (model.Token, error) {
	hash := hashToken(rawToken)
	token, err := s.authRepo.GetByHash(ctx, hash)
	if err != nil {
		return model.Token{}, fmt.Errorf("get refresh token: %w", err)
	}
	if token.Revoked {
		return model.Token{}, model.ErrTokenRevoked
	}
	if time.Now().After(token.ExpiresAt) {
		return model.Token{}, model.ErrTokenExpired
	}
	err = s.authRepo.Revoke(ctx, hash)
	if err != nil {
		return model.Token{}, fmt.Errorf("revoke refresh token: %w", err)
	}
	if err := s.authRepo.Revoke(ctx, hash); err != nil {
		return model.Token{}, fmt.Errorf("revoke refresh token: %w", err)
	}
	u, err := s.userRepo.GetUserByID(ctx, token.UserID)
	if err != nil {
		return model.Token{}, fmt.Errorf("get user by id: %w", err)
	}
	accessToken, err := s.newAccessToken(token.UserID, u.Role)
	if err != nil {
		return model.Token{}, fmt.Errorf("create access token: %w", err)
	}
	refreshToken, err := s.newRefreshToken(ctx, token.UserID)
	if err != nil {
		return model.Token{}, fmt.Errorf("create refresh token: %w", err)
	}
	return model.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
