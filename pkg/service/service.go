package service

import (
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	GetUserById(id int) (*data.User, error)
	UpdateUser(user data.User) error
	CreateUser(user data.User) (int, error)
	GetUser(email, password string) (*data.User, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos),
	}
}
