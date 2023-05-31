package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kirill-27/debt_manager/data"
)

type FriendsPostgres struct {
	db *sqlx.DB
}

func NewFriendsPostgres(db *sqlx.DB) *FriendsPostgres {
	return &FriendsPostgres{db: db}
}

func (f *FriendsPostgres) AddFriend(myId int, friendId int) error {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (my_id, friend_id)"+
		" values ($1, $2) RETURNING my_id", friendsTable)

	row := f.db.QueryRow(query, myId, friendId)
	err := row.Scan(&id)
	return err
}

func (f *FriendsPostgres) CheckIfFriendExists(myId int, friendId int) (bool, error) {
	var review data.Friends
	query := fmt.Sprintf("SELECT * FROM %s WHERE my_id=$1 AND friend_id=$2", friendsTable)
	err := f.db.Get(&review, query, myId, friendId)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, err
}
