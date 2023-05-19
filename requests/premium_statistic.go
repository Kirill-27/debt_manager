package requests

type PremiumStatistic struct {
	MyDebtsAmount      int `json:"my_debts_amount" binding:"required"`
	FriendsDebtsAmount int `json:"friends_debts_amount" binding:"required"`
	MyDebtsNumber      int `json:"my_debts_number" binding:"required"`
	FriendsDebtsNumber int `json:"friends_debts_number" binding:"required"`
}
