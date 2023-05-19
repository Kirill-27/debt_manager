package requests

type CommonStatistic struct {
	MyDebtsAmount         int                     `json:"my_debts_amount" binding:"required"`
	FriendsDebtsAmount    int                     `json:"friends_debts_amount" binding:"required"`
	MyDebtsNumber         int                     `json:"my_debts_number" binding:"required"`
	FriendsDebtsNumber    int                     `json:"friends_debts_number" binding:"required"`
	TopFriendsInteraction []TopFriendsInteraction `json:"top_friends_interaction" binding:"required"`
}

type TopFriendsInteraction struct {
	FriendId           int `json:"friend_id" binding:"required"`
	InteractionsNumber int `json:"interactions_number" binding:"required"`
}
