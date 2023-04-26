package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kirill-27/debt_manager/data"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user data.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (password, email, full_name, photo, rating, subscription_type)"+
		" values ($1, $2, $3, $4, $5, $6) RETURNING id", usersTable)

	row := r.db.QueryRow(query,
		user.Password,
		user.Email,
		user.FullName,
		user.Photo,
		user.Rating,
		user.SubscriptionType)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(email, password string) (*data.User, error) {
	var customer data.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password=$2", usersTable)
	err := r.db.Get(&customer, query, email, password)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &customer, err
}

func (r *AuthPostgres) GetUserById(id int) (*data.User, error) {
	var user data.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)

	return &user, err
}

func (r *AuthPostgres) UpdateUser(user data.User) error {
	//query := fmt.Sprintf("UPDATE %s SET name=$1, password=$2, username=$3, email=$4, phone=$5   WHERE id=$6 ", customersTable)
	//_, err := r.db.Exec(query, customer.Name, customer.Password, customer.Username, customer.Email, customer.Phone, customer.Id)
	return nil
}
