package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
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
