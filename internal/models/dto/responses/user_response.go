package responses

type UserData struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token,omitempty"`
}

type UserResponse struct {
	User UserData `json:"user"`
}
