package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kirill-27/debt_manager/data"
	"github.com/kirill-27/debt_manager/helpers"
	"strconv"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user data.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (password, email, full_name, photo, rating)"+
		" values ($1, $2, $3, $4, $5) RETURNING id", usersTable)

	row := r.db.QueryRow(query,
		user.Password,
		user.Email,
		user.FullName,
		user.Photo,
		user.Rating)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(email, password *string) (*data.User, error) {
	var user data.User
	var err error
	if email != nil && password != nil {
		query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password=$2", usersTable)
		err = r.db.Get(&user, query, email, password)
	} else {
		query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1", usersTable)
		err = r.db.Get(&user, query, email)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (r *AuthPostgres) GetAllUsers(sortBy []string, friendsFor *int) ([]data.User, error) {
	query := fmt.Sprintf("SELECT %s.* FROM %s ", usersTable, usersTable)

	var params []interface{}

	if friendsFor != nil {
		query += fmt.Sprintf("LEFT JOIN %s on %s.id =  %s.friend_id ", friendsTable, usersTable, friendsTable) +
			fmt.Sprintf("WHERE %s.my_id = $", friendsTable) + strconv.Itoa(len(params)+1)
		params = append(params, *friendsFor)
	}

	if len(sortBy) != 0 {
		query += helpers.ParseSortBy(sortBy)
	}

	rows, err := r.db.Queryx(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []data.User
	for rows.Next() {
		var user data.User
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *AuthPostgres) GetUserById(id int) (*data.User, error) {
	var user data.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (r *AuthPostgres) UpdateUser(user data.User) error {
	query := fmt.Sprintf(
		"UPDATE %s SET email=$2, password=$3, full_name=$4, subscription_type=$5, photo=$6, rating=$7 WHERE id=$1 ",
		usersTable)

	_, err := r.db.Exec(query, user.Id, user.Email, user.Password, user.FullName,
		user.SubscriptionType, user.Photo, user.Rating)
	return err
}
