package data

type CurrentDebts struct {
	Id       int `json:"id" db:"id"`
	DebtorID int `json:"debtor_id" db:"debtor_id" binding:"required"`
	LenderId int `json:"lender_id" db:"lender_id" binding:"required"`
	Amount   int `json:"amount" db:"amount" binding:"required"`
}
