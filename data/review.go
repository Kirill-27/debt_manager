package data

import "time"

type Review struct {
	Id         int       `json:"id" db:"id"`
	ReviewerId int       `json:"reviewer_id" db:"reviewer_id"`
	LenderId   int       `json:"lender_id" db:"lender_id" binding:"required"`
	Comment    string    `json:"comment" db:"comment" binding:"required"`
	Rate       int       `json:"rate" db:"rate" binding:"required"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
