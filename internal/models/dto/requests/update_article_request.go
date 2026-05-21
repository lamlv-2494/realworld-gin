package requests

type UpdateArticleData struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Body        *string `json:"body"`
}

type UpdateArticleRequest struct {
	Article UpdateArticleData `json:"article"`
}
