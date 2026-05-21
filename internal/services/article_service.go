package services

import (
	"realworld-gin/internal/models"
	"realworld-gin/internal/models/dto/responses"
	"realworld-gin/internal/repositories"
	"realworld-gin/internal/utils/constants"
	"realworld-gin/internal/utils/slug_gen"
)

type ArticleService interface {
	CreateArticle(authorID uint, title, description, body string, tagList []string) (*responses.ArticleResponse, error)

	GetArticle(currentUserID uint, slug string) (*responses.ArticleResponse, error)

	UpdateArticle(currentUserID uint, slug string, data map[string]any) (*responses.ArticleResponse, error)

	DeleteArticle(currentUserID uint, slug string) error

	ListArticles(currentUserID uint, tag, author, favorited string, limit, page int) (*responses.ArticleListResponse, error)

	FeedArticle(currentUserID uint, limit, page int) (*responses.ArticleListResponse, error)

	GetTags() (*responses.TagsResponse, error)
}

type articleService struct {
	articleRepo  repositories.ArticleRepository
	userRepo     repositories.UserRepository
	favoriteRepo repositories.FavoriteRepository
}

func NewArticleService(articleRepo repositories.ArticleRepository, userRepo repositories.UserRepository, favoriteRepo repositories.FavoriteRepository) ArticleService {
	return &articleService{
		articleRepo:  articleRepo,
		userRepo:     userRepo,
		favoriteRepo: favoriteRepo,
	}
}

// CreateArticle implements [ArticleService].
func (a *articleService) CreateArticle(authorID uint, title string, description string, body string, tagList []string) (*responses.ArticleResponse, error) {
	author, err := a.userRepo.FindByID(authorID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	slug := slug_gen.GenerateSlug(title)

	var tags []*models.Tag
	for _, tagName := range tagList {
		tag, err := a.articleRepo.FindOrCreateTag(tagName)
		if err == nil {
			tags = append(tags, tag)
		}
	}

	article := &models.Article{
		Slug:        slug,
		Title:       title,
		Description: description,
		Body:        body,
		AuthorID:    authorID,
		Tags:        tags,
	}

	if err := a.articleRepo.CreateArticle(article); err != nil {
		return nil, ErrCreateArticle
	}

	response := &responses.ArticleResponse{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        tagList,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		Favorited:      false, // Mặc định bài mới tạo chưa có ai like
		FavoritesCount: 0,
		Author: responses.ProfileResponse{
			Profile: responses.ProfileData{
				Username:  author.Username,
				Bio:       author.Bio,
				Image:     author.Image,
				Following: false, // Tự mình không follow mình
			},
		},
	}

	return response, nil
}

// GetArticle implements [ArticleService].
func (a *articleService) GetArticle(currentUserID uint, slug string) (*responses.ArticleResponse, error) {
	article, err := a.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, ErrArticleNotFound
	}

	var tagList []string
	for _, tag := range article.Tags {
		tagList = append(tagList, tag.Name)
	}

	isFollowing := false
	if currentUserID != 0 {
		isFollowing = a.userRepo.IsFollowing(currentUserID, article.AuthorID)
	}

	isFavorited := false
	if currentUserID != 0 {
		isFavorited = a.favoriteRepo.IsFavorited(currentUserID, article.ID)
	}
	favoritesCount := a.favoriteRepo.GetFavoritesCount(article.ID)

	response := responses.ArticleResponse{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        tagList,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		Favorited:      isFavorited,
		FavoritesCount: favoritesCount,
		Author: responses.ProfileResponse{
			Profile: responses.ProfileData{
				Username:  article.Author.Username,
				Bio:       article.Author.Bio,
				Image:     article.Author.Image,
				Following: isFollowing,
			},
		},
	}

	return &response, nil

}

// UpdateArticle implements [ArticleService].
func (a *articleService) UpdateArticle(currentUserID uint, slug string, data map[string]any) (*responses.ArticleResponse, error) {
	article, err := a.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, ErrArticleNotFound
	}

	if article.AuthorID != currentUserID {
		return nil, ErrNotArticleAuthor
	}

	if title, ok := data[constants.Title]; ok {
		newSlug := slug_gen.GenerateSlug(title.(string))
		data[constants.Slug] = newSlug
	}

	if err := a.articleRepo.UpdateArticle(article, data); err != nil {
		return nil, ErrNotArticleAuthor
	}

	finalSlug := article.Slug
	if newSlug, ok := data[constants.Slug]; ok {
		finalSlug = newSlug.(string)
	}

	return a.GetArticle(currentUserID, finalSlug)
}

// DeleteArticle implements [ArticleService].
func (a *articleService) DeleteArticle(currentUserID uint, slug string) error {
	article, err := a.articleRepo.FindBySlug(slug)
	if err != nil {
		return ErrArticleNotFound
	}

	if article.AuthorID != currentUserID {
		return ErrNotArticleAuthor
	}

	if err := a.articleRepo.Delete(article); err != nil {
		return ErrDeleteArticle
	}
	return nil
}

// ListArticles implements [ArticleService].
func (a *articleService) ListArticles(currentUserID uint, tag, author, favorited string, limit, page int) (*responses.ArticleListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	articles, totalCount, err := a.articleRepo.ListArticles(tag, author, favorited, limit, offset)
	if err != nil {
		return nil, ErrArticleNotFound
	}

	var articlesResponse []responses.ArticleResponse

	if len(articles) == 0 {
		articlesResponse = make([]responses.ArticleResponse, 0)
	}

	for _, article := range articles {
		var tagList []string
		for _, tag := range article.Tags {
			tagList = append(tagList, tag.Name)
		}

		isFollowing := false
		if currentUserID != 0 {
			isFollowing = a.userRepo.IsFollowing(currentUserID, article.AuthorID)
		}

		articlesResponse = append(articlesResponse, responses.ArticleResponse{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        tagList,
			CreatedAt:      article.CreatedAt,
			UpdatedAt:      article.UpdatedAt,
			Favorited:      false,
			FavoritesCount: 0,
			Author: responses.ProfileResponse{
				Profile: responses.ProfileData{
					Username:  article.Author.Username,
					Bio:       article.Author.Bio,
					Image:     article.Author.Image,
					Following: isFollowing,
				},
			},
		})
	}
	return &responses.ArticleListResponse{
		Articles:   articlesResponse,
		TotalCount: totalCount,
	}, nil
}

// FeedArticle implements [ArticleService].
func (a *articleService) FeedArticle(currentUserID uint, limit int, page int) (*responses.ArticleListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	articles, totalCount, err := a.articleRepo.FeedArticle(currentUserID, limit, offset)
	if err != nil {
		return nil, err
	}

	var articlesResponse []responses.ArticleResponse
	if len(articles) == 0 {
		articlesResponse = make([]responses.ArticleResponse, 0)
	}

	for _, article := range articles {
		var tagList []string
		for _, tag := range article.Tags {
			tagList = append(tagList, tag.Name)
		}

		articlesResponse = append(articlesResponse, responses.ArticleResponse{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        tagList,
			CreatedAt:      article.CreatedAt,
			UpdatedAt:      article.UpdatedAt,
			Favorited:      false, // Sẽ làm ở Module Favorites
			FavoritesCount: 0,
			Author: responses.ProfileResponse{
				Profile: responses.ProfileData{
					Username:  article.Author.Username,
					Bio:       article.Author.Bio,
					Image:     article.Author.Image,
					Following: true,
				},
			},
		})
	}

	response := responses.ArticleListResponse{
		Articles:   articlesResponse,
		TotalCount: totalCount,
	}

	return &response, nil
}

func (s *articleService) GetTags() (*responses.TagsResponse, error) {
	tag, err := s.articleRepo.GetTags()
	if err != nil {
		return nil, ErrTagsNotFound
	}

	res := responses.TagsResponse{
		Tags: tag,
	}

	return &res, nil
}
