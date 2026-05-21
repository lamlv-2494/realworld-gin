package requests

type CommentRequest struct {
	Comment struct {
		Body string `json:"body" binding:"required"`
	} `json:"comment"`
}
