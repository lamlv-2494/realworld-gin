package responses

type ProfileData struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type ProfileResponse struct {
	Profile ProfileData `json:"profile"`
}
