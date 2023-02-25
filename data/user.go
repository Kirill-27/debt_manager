package data

const (
	SubscriptionTypeFree = iota + 1
	SubscriptionTypePremium
)

type User struct {
	Id               int     `json:"id" db:"id"`
	Email            string  `json:"email" db:"email" binding:"required"`
	Password         string  `json:"password" db:"password" binding:"required"`
	FullName         string  `json:"full_name" db:"full_name" binding:"required"`
	SubscriptionType int     `json:"subscription_type" db:"subscription_type"`
	Photo            string  `json:"photo" db:"photo"`
	Rating           float64 `json:"rating" db:"rating"`
}
