package service

import (
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	GetAllUsers(sortBy []string, friendsFor *int) ([]data.User, error)
	GetUserById(id int) (*data.User, error)
	UpdateUser(user data.User) error
	CreateUser(user data.User) (int, error)
	GetUser(email, password *string) (*data.User, error)
	GenerateToken(email, password string) (int, string, error)
	ParseToken(token string) (int, error)
}

type Debt interface {
	CreateDebt(debt data.Debt) (int, error)
	GetAllDebts(debtorId *int, lenderId *int, statuses string, sortBy []string) ([]data.Debt, error)
	GetDebtById(debtId int) (*data.Debt, error)
	UpdateStatus(id int, status int) error
	DeleteDebt(debtId int) error
	CloseAllDebts(debtorId int, lenderId int) error
}

type CurrentDebt interface {
	CreateCurrentDebt(debt data.CurrentDebts) (int, error)
	GetAllCurrentDebts(debtorId *int, lenderId *int, sortBy []string) ([]data.CurrentDebts, error)
	UpdateAmount(id int, amount int) error
	DeleteCurrentDebt(debtId int) error
}

type Review interface {
	GetReviewById(id int) (*data.Review, error)
	GetAllReviews(reviewerId *int, lenderId *int, sortBy []string) ([]data.Review, error)
	UpdateReview(review data.Review) error
	CreateReview(review data.Review) (int, error)
}

type Friends interface {
	AddFriend(myId int, friendId int) error
	CheckIfFriendExists(myId int, friendId int) (bool, error)
}

type Service struct {
	Authorization
	Debt
	CurrentDebt
	Review
	Friends
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos),
		Debt:          NewDebtService(repos),
		CurrentDebt:   NewCurrentDebtService(repos),
		Review:        NewReviewService(repos),
		Friends:       NewFriendsService(repos),
	}
}
