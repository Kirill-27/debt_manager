package data

type Friends struct {
	MyId     int `json:"my_id" db:"my_id" binding:"required"`
	FriendId int `json:"friend_id" db:"friend_id" binding:"required"`
}
