package data

import "time"

type Review struct {
	Id        int       `json:"id" db:"id"`
	DebtorID  int       `json:"debtor_id" db:"debtor_id" binding:"required"`
	LenderId  int       `json:"lender_id" db:"lender_id" binding:"required"`
	Comment   string    `json:"comment" db:"comment" binding:"required"`
	Rate      int       `json:"rate" db:"rate" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at" binding:"required"`
}
