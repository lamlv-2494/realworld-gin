package requests

type UpdateUserRequest struct {
	User struct {
		Email    *string `json:"email"`
		Username *string `json:"username"`
		Password *string `json:"password"`
		Image    *string `json:"image"`
		Bio      *string `json:"bio"`
	} `json:"user"`
}
