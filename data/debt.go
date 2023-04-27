package data

import "time"

const (
	DebtStatusPendingActive = iota + 1
	DebtStatusActive
	DebtStatusClosed
)

type Debt struct {
	Id          int       `json:"id" db:"id"`
	DebtorId    int       `json:"debtor_id" db:"debtor_id" binding:"required"`
	LenderId    int       `json:"lender_id" db:"lender_id" binding:"required"`
	Status      int       `json:"status" db:"status" `
	Amount      int       `json:"amount" db:"amount" binding:"required"`
	Description string    `json:"description" db:"description" binding:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" `
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" `
}
