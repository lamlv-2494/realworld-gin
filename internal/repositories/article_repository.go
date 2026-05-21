package repositories

import (
	config "realworld-gin/internal/config/database"
	"realworld-gin/internal/models"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	CreateArticle(article *models.Article) error

	FindOrCreateTag(name string) (*models.Tag, error)

	FindBySlug(slug string) (*models.Article, error)

	UpdateArticle(article *models.Article, data map[string]any) error

	Delete(article *models.Article) error

	ListArticles(tag, author, favorited string, limit, offset int) ([]*models.Article, int64, error)

	FeedArticle(followerID uint, limit, offset int) ([]*models.Article, int64, error)

	GetTags() ([]string, error)
}

type articleRepository struct {
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{}
}

// CreateArticle implements [ArticleRepository].
func (a *articleRepository) CreateArticle(article *models.Article) error {
	return config.DB.
		Create(article).
		Error
}

// GetArticle implements [ArticleRepository].
func (a *articleRepository) FindOrCreateTag(name string) (*models.Tag, error) {
	var tag models.Tag

	err := config.DB.
		Where("name = ?", name).
		FirstOrCreate(&tag, models.Tag{Name: name}).
		Error

	return &tag, err
}

// GetArticle implements [ArticleRepository].
func (a *articleRepository) FindBySlug(slug string) (*models.Article, error) {
	var article models.Article
	err := config.DB.
		Preload("Author").
		Preload("Tags").
		Where("slug = ?", slug).
		First(&article).
		Error

	if err != nil {
		return nil, err
	}
	return &article, nil
}

// UpdateArticle implements [ArticleRepository].
func (a *articleRepository) UpdateArticle(article *models.Article, data map[string]any) error {
	return config.DB.
		Model(article).
		Updates(data).
		Error
}

// Delete implements [ArticleRepository].
func (a *articleRepository) Delete(article *models.Article) error {
	return config.DB.
		Delete(article).
		Error
}

// ListArticles implements [ArticleRepository].
func (a *articleRepository) ListArticles(tag string, author string, favorited string, limit int, offset int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var count int64

	query := config.DB.Model(&models.Article{})
	if author != "" {
		query = query.Where("author_id IN (SELECT id FROM users WHERE username = ?)", author)
	}
	if tag != "" {
		query = query.Where("id IN (SELECT article_id FROM article_tags INNER JOIN tags ON tags.id = article_tags.tag_id WHERE tags.name = ?)", tag)
	}

	query.Count(&count)

	err := query.
		Preload("Author").
		Preload("Tags").
		Order("articles.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).
		Error

	return articles, count, err
}

// FeedArticle implements [ArticleRepository].
func (a *articleRepository) FeedArticle(followerID uint, limit int, offset int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var totalCount int64

	query := config.DB.Model(&models.Article{}).
		Where("author_id IN (SELECT following_id FROM user_follows WHERE follower_id = ?)", followerID)

	query.Count(&totalCount)

	err := query.
		Preload("Author").
		Preload("Tags").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).
		Error

	return articles, totalCount, err
}

// GetTags implements [ArticleRepository].
func (a *articleRepository) GetTags() ([]string, error) {
	var tagName []string

	err := config.DB.
		Model(&models.Tag{}).
		Pluck("name", &tagName).
		Error

	return tagName, err
}
