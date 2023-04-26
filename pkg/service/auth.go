package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/pkg/repository"
	"time"
)

const (
	signingKey = "dfvjdfv@21#@d(q*Djdsdf"
	tokenTTL   = 60 * 24 * time.Hour
	salt       = "h0q12hqw124f17ajf3ajs"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user data.User) (int, error) {
	user.SubscriptionType = data.SubscriptionTypeFree
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(email, password string) (int, string, error) {
	user, err := s.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return 0, "", err
	}
	if user == nil {
		return 0, "", errors.New("wrong username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	tokenString, err := token.SignedString([]byte(signingKey))
	return user.Id, tokenString, err
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) GetUserById(id int) (*data.User, error) {
	return s.repo.GetUserById(id)
}

func (s *AuthService) UpdateUser(user data.User) error {
	return s.repo.UpdateUser(user)
}

func (s *AuthService) GetUser(email, password string) (*data.User, error) {
	return s.repo.GetUser(email, password)
}
