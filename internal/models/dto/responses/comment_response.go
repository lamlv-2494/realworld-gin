package responses

import "time"

type CommentResponse struct {
	ID        uint            `json:"id"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Body      string          `json:"body"`
	Author    ProfileResponse `json:"author"`
}
