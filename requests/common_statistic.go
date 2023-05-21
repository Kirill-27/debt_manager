package requests

type CommonStatistic struct {
	MyActiveDebtsAmount      int `json:"my_active_debts_amount" binding:"required"`
	FriendsActiveDebtsAmount int `json:"friends_active_debts_amount" binding:"required"`
	MyActiveDebtsNumber      int `json:"my_active_debts_number" binding:"required"`
	FriendsActiveDebtsNumber int `json:"friends_active_debts_number" binding:"required"`
	TopLender                int `json:"top_lender" binding:"required"`
	TopDebtor                int `json:"top_debtor" binding:"required"`
}
