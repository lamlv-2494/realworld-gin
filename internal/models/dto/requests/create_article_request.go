package requests

type ArticleData struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Body        string   `json:"body" binding:"required"`
	TagList     []string `json:"tagList"`
}

type CreateArticleRequest struct {
	Article ArticleData `json:"article"`
}
