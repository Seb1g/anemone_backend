package auth_services

import (
	"anemone_notes/internal/model/auth_model"
	"anemone_notes/internal/repository/auth_repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var accessKey = []byte("access-secret")
var refreshKey = []byte("refresh-secret")

type AuthService struct {
	Users   *auth_repository.UserRepo
	Refresh *auth_repository.RefreshRepo
}

func NewAuthService(u *auth_repository.UserRepo, r *auth_repository.RefreshRepo) *AuthService {
	return &AuthService{Users: u, Refresh: r}
}

// ─── регистрация ──────────────────────────────────────────────
func (s *AuthService) Register(ctx context.Context, email, password string) (*auth_model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := &auth_model.User{Email: email, Password: string(hash)}
	if err := s.Users.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

// ─── генерация токенов ─────────────────────────────────────────
func (s *AuthService) generateTokens(ctx context.Context, u *auth_model.User) (string, string, error) {
	// access (15 минут)
	accessClaims := jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(60 * time.Minute).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := at.SignedString(accessKey)
	if err != nil {
		return "", "", err
	}

	// refresh (7 дней)
	refreshClaims := jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := rt.SignedString(refreshKey)
	if err != nil {
		return "", "", err
	}

	// сохранить refresh в БД
	if err := s.Refresh.Store(ctx, u.ID, refreshToken, time.Now().Add(7*24*time.Hour)); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ─── логин ─────────────────────────────────────────────────────
func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	u, err := s.Users.GetByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("user not found")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return "", "", errors.New("invalid credentials")
	}

	return s.generateTokens(ctx, u)
}

// ─── сброс пароля ──────────────────────────────────────────────
func (s *AuthService) ResetPassword(ctx context.Context, email, newPassword string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	return s.Users.UpdatePassword(ctx, email, string(hash))
}

// ─── обновление access по refresh ─────────────────────────────
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshKey, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	userID := int(claims["user_id"].(float64))
	exp := int64(claims["exp"].(float64))

	ok, err := s.Refresh.Check(ctx, userID, refreshToken, time.Unix(exp, 0))
	if err != nil || !ok {
		return "", errors.New("refresh token not found or expired")
	}

	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	return at.SignedString(accessKey)
}

func (s *AuthService) ParseAccessToken(tokenStr string) (int, error) {
    claims := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return accessKey, nil
    })
    if err != nil || !token.Valid {
        return 0, errors.New("invalid token")
    }
    userID := int(claims["user_id"].(float64))
    return userID, nil
}