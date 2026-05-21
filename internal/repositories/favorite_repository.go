package repositories

import (
	config "realworld-gin/internal/config/database"
	"realworld-gin/internal/models"

	"gorm.io/gorm"
)

type FavoriteRepository interface {
	FavoriteArticle(user *models.User, article *models.Article) error
	UnfavoriteArticle(user *models.User, article *models.Article) error
	IsFavorited(userID, articleID uint) bool
	GetFavoritesCount(articleID uint) int64
}

type favoriteRepository struct {
}

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{}
}

// FavoriteArticle implements [FavoriteRepository].
func (f *favoriteRepository) FavoriteArticle(user *models.User, article *models.Article) error {
	return config.DB.
		Model(user).
		Association("Favorites").
		Append(article)
}

// UnfavoriteArticle implements [FavoriteRepository].
func (f *favoriteRepository) UnfavoriteArticle(user *models.User, article *models.Article) error {
	return config.DB.
		Model(user).
		Association("Favorites").
		Delete(article)
}

// GetFavoritesCount implements [FavoriteRepository].
func (f *favoriteRepository) GetFavoritesCount(articleID uint) int64 {
	var count int64

	config.DB.
		Table("user_favorites").
		Where("article_id = ?", articleID).
		Count(&count)

	return count
}

// IsFavorited implements [FavoriteRepository].
func (f *favoriteRepository) IsFavorited(userID uint, articleID uint) bool {
	var count int64

	config.DB.
		Table("user_favorites").
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count)

	return count > 0
}
