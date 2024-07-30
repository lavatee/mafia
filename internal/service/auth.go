package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lavatee/mafia/internal/repository"
)

const (
	salt       = "jhakldfhaae"
	accessTTL  = 20 * time.Second
	refreshTTL = 40 * time.Second
	tokenKey   = "qeq0efquj"
)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) HashPassword(password string) string {
	sha := sha1.New()
	sha.Write([]byte(password))
	return fmt.Sprintf("%x", sha.Sum([]byte(salt)))
}

func (s *AuthService) SignUp(email string, name string, password string) (int, error) {
	return s.repo.SignUp(email, name, s.HashPassword(password))
}

func (s *AuthService) SignIn(email string, password string) (string, string, error) {
	user, err := s.repo.SignIn(email, s.HashPassword(password))
	if err != nil {
		return "", "", err
	}
	accessClaims := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(accessTTL).Unix(),
	}
	refreshClaims := jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(accessTTL).Unix(),
	}
	access, err := s.NewToken(accessClaims)
	if err != nil {
		return "", "", err
	}
	refresh, err := s.NewToken(refreshClaims)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil

}
func (s *AuthService) NewToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *AuthService) Refresh(token string) (string, string, error) {
	parsedToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(tokenKey), nil
	})
	if err != nil {
		return "", "", err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		access, err := s.NewToken(jwt.MapClaims{
			"id":  claims["id"],
			"exp": time.Now().Add(accessTTL).Unix(),
		})
		if err != nil {
			return "", "", err
		}
		refresh, err := s.NewToken(jwt.MapClaims{
			"id":  claims["id"],
			"exp": time.Now().Add(refreshTTL).Unix(),
		})
		if err != nil {
			return "", "", err
		}
		return access, refresh, nil
	}
	return "", "", errors.New("token is expired falos")
}
