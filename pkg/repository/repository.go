package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kirill-27/debt_manager/data"
)

type Authorization interface {
	GetUserById(id int) (*data.User, error)
	UpdateUser(user data.User) error
	CreateUser(user data.User) (int, error)
	GetUser(email, password *string) (*data.User, error)
	GetAllUsers(sortBy []string, friendsFor *int) ([]data.User, error)
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

type Repository struct {
	Authorization
	Debt
	CurrentDebt
	Review
	Friends
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Debt:          NewDebtPostgres(db),
		CurrentDebt:   NewCurrentDebtPostgres(db),
		Review:        NewReviewPostgres(db),
		Friends:       NewFriendsPostgres(db),
	}
}
