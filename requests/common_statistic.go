package requests

type CommonStatistic struct {
	MyActiveDebtsAmount      int   `json:"my_active_debts_amount" binding:"required"`
	FriendsActiveDebtsAmount int   `json:"friends_active_debts_amount" binding:"required"`
	MyActiveDebtsNumber      int   `json:"my_active_debts_number" binding:"required"`
	FriendsActiveDebtsNumber int   `json:"friends_active_debts_number" binding:"required"`
	TopLenders               []int `json:"top_lenders" binding:"required"`
	TopDebtors               []int `json:"top_debtors" binding:"required"`
}
