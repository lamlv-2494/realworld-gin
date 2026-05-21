package responses

import "time"

type ArticleResponse struct {
	Slug           string          `json:"slug"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Body           string          `json:"body"`
	TagList        []string        `json:"tagList"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	Favorited      bool            `json:"favorited"`
	FavoritesCount int64           `json:"favoritesCount"`
	Author         ProfileResponse `json:"author"`
}

type ArticleListResponse struct {
	Articles   []ArticleResponse `json:"articles"`
	TotalCount int64             `json:"total_count"`
}
