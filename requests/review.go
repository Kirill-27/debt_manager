package requests

type UpdateReview struct {
	Comment string `json:"comment"`
	Rate    int    `json:"rate"`
}
