package repositories

import (
	config "realworld-gin/internal/config/database"
	"realworld-gin/internal/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(comment *models.Comment) error
	GetCommentsByArticleID(articleID uint) ([]*models.Comment, error)
	GetCommentByID(id uint) (*models.Comment, error)
	DeleteComment(comment *models.Comment) error
}

type commentRepository struct{}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{}
}

// CreateComment implements [CommentRepository].
func (c *commentRepository) CreateComment(comment *models.Comment) error {
	return config.DB.
		Create(comment).
		Error
}

// GetCommentByID implements [CommentRepository].
func (c *commentRepository) GetCommentByID(id uint) (*models.Comment, error) {
	var comment models.Comment

	err := config.DB.
		First(&comment, id).
		Error

	return &comment, err
}

// GetCommentsByArticleID implements [CommentRepository].
func (c *commentRepository) GetCommentsByArticleID(articleID uint) ([]*models.Comment, error) {

	var comments []*models.Comment

	err := config.DB.
		Preload("Author").
		Where("article_id = ?", articleID).
		Order("created_at DESC").
		Find(&comments).
		Error

	return comments, err
}

// DeleteComment implements [CommentRepository].
func (c *commentRepository) DeleteComment(comment *models.Comment) error {
	return config.DB.
		Delete(comment).
		Error
}
