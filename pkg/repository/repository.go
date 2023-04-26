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
}

type Debt interface {
	CreateDebt(debt data.Debt) (int, error)
	GetAllDebts(debtorId *int, lenderId *int, sortBy []string) ([]data.Debt, error)
	GetDebtById(debtId int) (*data.Debt, error)
	UpdateDebt(debt data.Debt) error
	DeleteDebt(debtId int) error
}

type Repository struct {
	Authorization
	Debt
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Debt:          NewDebtPostgres(db),
	}
}
