package services

import (
	"realworld-gin/internal/models"
	"realworld-gin/internal/models/dto/responses"
	"realworld-gin/internal/repositories"
)

type CommentService interface {
	CreateComment(currentUserID uint, slug string, body string) (*responses.CommentResponse, error)
	GetComments(currentUserID uint, slug string) ([]responses.CommentResponse, error)
	DeleteComment(currentUserID uint, slug string, commentID uint) error
}

type commentService struct {
	commentRepo repositories.CommentRepository
	articleRepo repositories.ArticleRepository
	userRepo    repositories.UserRepository
}

func NewCommentService(c repositories.CommentRepository, a repositories.ArticleRepository, u repositories.UserRepository) CommentService {
	return &commentService{c, a, u}
}

// CreateComment implements [CommentService].
func (s *commentService) CreateComment(currentUserID uint, slug string, body string) (*responses.CommentResponse, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, ErrArticleNotFound
	}

	comment := &models.Comment{
		Body:      body,
		ArticleID: article.ID,
		AuthorID:  currentUserID,
	}

	if err := s.commentRepo.CreateComment(comment); err != nil {
		return nil, ErrCreateComment
	}

	author, err := s.userRepo.FindByID(currentUserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	commentResponse := &responses.CommentResponse{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		Body:      comment.Body,
		Author: responses.ProfileResponse{
			Profile: responses.ProfileData{
				Username:  author.Username,
				Bio:       author.Bio,
				Image:     author.Image,
				Following: false,
			},
		},
	}

	return commentResponse, nil
}

// DeleteComment implements [CommentService].
func (s *commentService) DeleteComment(currentUserID uint, slug string, commentID uint) error {
	_, err := s.userRepo.FindByID(currentUserID)
	if err != nil {
		return ErrUserNotFound
	}

	comment, err := s.commentRepo.GetCommentByID(commentID)
	if err != nil {
		return ErrCommentNotFound
	}

	err = s.commentRepo.DeleteComment(comment)
	if err != nil {
		return ErrNotCommentAuthor
	}
	return nil

}

// GetComments implements [CommentService].
func (s *commentService) GetComments(currentUserID uint, slug string) ([]responses.CommentResponse, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, ErrArticleNotFound
	}

	comments, err := s.commentRepo.GetCommentsByArticleID(article.ID)
	if err != nil {
		return nil, ErrCommentNotFound
	}

	var commentsResponse = make([]responses.CommentResponse, 0)
	for _, c := range comments {
		isFollowing := false

		if currentUserID != 0 {
			isFollowing = s.userRepo.IsFollowing(currentUserID, c.AuthorID)
		}

		commentsResponse = append(commentsResponse, responses.CommentResponse{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			Author: responses.ProfileResponse{Profile: responses.ProfileData{
				Username:  c.Author.Username,
				Bio:       c.Author.Bio,
				Image:     c.Author.Image,
				Following: isFollowing,
			}},
		})
	}
	return commentsResponse, nil

}
