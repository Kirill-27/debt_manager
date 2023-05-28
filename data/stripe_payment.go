package data

import "time"

const (
	StripePaymentsStatusProcessing = iota + 1
	StripePaymentsStatusSucceeded
	StripePaymentsStatusCanceled
	StripePaymentsStatusIncorrectEmail
)

type StripePayment struct {
	PaymentId string    `json:"payment_id" db:"payment_id" binding:"required"`
	UserId    int       `json:"user_id" db:"user_id" binding:"required"`
	Status    int       `json:"status" db:"status" `
	CreatedAt time.Time `json:"created_at" db:"created_at" `
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" `
}
