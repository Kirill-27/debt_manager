package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/pkg/repository"
	"time"
)

const (
	signingKey = "dfvjdfv@21#@d(q*Djdsdf"
	tokenTTL   = 60 * 24 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	CustomerId int `json:"customer_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(customer data.User) (int, error) {
	return s.repo.CreateUser(customer)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	customer, err := s.repo.GetUser(username, password)
	if err != nil {
		return "", err
	}
	if customer == nil {
		return "", errors.New("wrong username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		customer.Id,
	})

	return token.SignedString([]byte(signingKey))
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

	return claims.CustomerId, nil
}

func (s *AuthService) GetUserById(id int) (*data.User, error) {
	return s.repo.GetUserById(id)
}

func (s *AuthService) UpdateUser(user data.User) error {
	return s.repo.UpdateUser(user)
}

func (s *AuthService) GetUser(username, password string) (*data.User, error) {
	return s.repo.GetUser(username, password)
}
