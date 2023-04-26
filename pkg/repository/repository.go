package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kirill-27/debt_manager/data"
)

type Authorization interface {
	GetUserById(id int) (*data.User, error)
	UpdateUser(user data.User) error
	CreateUser(user data.User) (int, error)
	GetUser(email, password string) (*data.User, error)
	GetAllUsers(sortBy []string) ([]data.User, error)
}

type Debt interface {
	CreateDebt(debt data.Debt) (int, error)
	GetAllDebts(debtorId *int, lenderId *int, sortBy []string) ([]data.Debt, error)
	GetDebtById(debtId int) (*data.Debt, error)
	UpdateStatus(id int, status int) error
	DeleteDebt(debtId int) error
}

type CurrentDebt interface {
	CreateCurrentDebt(debt data.CurrentDebts) (int, error)
	GetAllCurrentDebts(debtorId *int, lenderId *int, sortBy []string) ([]data.CurrentDebts, error)
	UpdateAmount(id int, amount int) error
	DeleteCurrentDebt(debtId int) error
}

type Repository struct {
	Authorization
	Debt
	CurrentDebt
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Debt:          NewDebtPostgres(db),
		CurrentDebt:   NewCurrentDebtPostgres(db),
	}
}
