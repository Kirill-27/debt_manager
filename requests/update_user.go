package requests

type UpdateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Photo    string `json:"photo"`
}
