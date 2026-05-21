package services

import (
	"realworld-gin/internal/repositories"
)

type FavoriteService interface {
	Favorite(currentUserID uint, slug string) error
	Unfavorite(currentUserID uint, slug string) error
}

type favoriteService struct {
	favoriteRepo repositories.FavoriteRepository
	userRepo     repositories.UserRepository
	articleRepo  repositories.ArticleRepository
}

func NewFavoriteService(
	favoriteRepo repositories.FavoriteRepository,
	userRepo repositories.UserRepository,
	articleRepo repositories.ArticleRepository,
) FavoriteService {
	return &favoriteService{favoriteRepo, userRepo, articleRepo}
}

// Favorite implements [FavoriteService].
func (f *favoriteService) Favorite(currentUserID uint, slug string) error {
	user, err := f.userRepo.FindByID(currentUserID)
	if err != nil {
		return err
	}

	article, err := f.articleRepo.FindBySlug(slug)
	if err != nil {
		return ErrArticleNotFound
	}
	return f.favoriteRepo.FavoriteArticle(user, article)
}

// Unfavorite implements [FavoriteService].
func (f *favoriteService) Unfavorite(currentUserID uint, slug string) error {
	user, err := f.userRepo.FindByID(currentUserID)
	if err != nil {
		return ErrUserNotFound
	}

	article, err := f.articleRepo.FindBySlug(slug)
	if err != nil {
		return ErrArticleNotFound
	}

	return f.favoriteRepo.UnfavoriteArticle(user, article)

}
