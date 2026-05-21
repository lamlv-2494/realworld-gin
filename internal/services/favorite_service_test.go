package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFavoriteService_Favorite_Success(t *testing.T) {
	favRepo := new(MockFavoriteRepository)
	userRepo := new(MockUserRepository)
	articleRepo := new(MockArticleRepository)
	svc := NewFavoriteService(favRepo, userRepo, articleRepo)

	user := stubUser(1, "alice", "alice@example.com")
	article := stubArticle(10, "cool-article", "Cool Article", 5)

	userRepo.On("FindByID", uint(1)).Return(user, nil)
	articleRepo.On("FindBySlug", "cool-article").Return(article, nil)
	favRepo.On("FavoriteArticle", user, article).Return(nil)

	err := svc.Favorite(1, "cool-article")
	require.NoError(t, err)
	favRepo.AssertExpectations(t)
}

func TestFavoriteService_Favorite_UserNotFound(t *testing.T) {
	favRepo := new(MockFavoriteRepository)
	userRepo := new(MockUserRepository)
	articleRepo := new(MockArticleRepository)
	svc := NewFavoriteService(favRepo, userRepo, articleRepo)

	userRepo.On("FindByID", uint(99)).Return(nil, ErrUserNotFound)

	err := svc.Favorite(99, "some-article")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestFavoriteService_Favorite_ArticleNotFound(t *testing.T) {
	favRepo := new(MockFavoriteRepository)
	userRepo := new(MockUserRepository)
	articleRepo := new(MockArticleRepository)
	svc := NewFavoriteService(favRepo, userRepo, articleRepo)

	user := stubUser(1, "alice", "alice@example.com")
	userRepo.On("FindByID", uint(1)).Return(user, nil)
	articleRepo.On("FindBySlug", "ghost").Return(nil, ErrArticleNotFound)

	err := svc.Favorite(1, "ghost")
	assert.ErrorIs(t, err, ErrArticleNotFound)
}

// ---------------------------------------------------------------------------
// Unfavorite
// ---------------------------------------------------------------------------

func TestFavoriteService_Unfavorite_Success(t *testing.T) {
	favRepo := new(MockFavoriteRepository)
	userRepo := new(MockUserRepository)
	articleRepo := new(MockArticleRepository)
	svc := NewFavoriteService(favRepo, userRepo, articleRepo)

	user := stubUser(1, "alice", "alice@example.com")
	article := stubArticle(10, "cool-article", "Cool Article", 5)

	userRepo.On("FindByID", uint(1)).Return(user, nil)
	articleRepo.On("FindBySlug", "cool-article").Return(article, nil)
	favRepo.On("UnfavoriteArticle", user, article).Return(nil)

	err := svc.Unfavorite(1, "cool-article")
	require.NoError(t, err)
	favRepo.AssertExpectations(t)
}

func TestFavoriteService_Unfavorite_UserNotFound(t *testing.T) {
	favRepo := new(MockFavoriteRepository)
	userRepo := new(MockUserRepository)
	articleRepo := new(MockArticleRepository)
	svc := NewFavoriteService(favRepo, userRepo, articleRepo)

	userRepo.On("FindByID", uint(99)).Return(nil, ErrUserNotFound)

	err := svc.Unfavorite(99, "some-article")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestFavoriteService_Unfavorite_ArticleNotFound(t *testing.T) {
	favRepo := new(MockFavoriteRepository)
	userRepo := new(MockUserRepository)
	articleRepo := new(MockArticleRepository)
	svc := NewFavoriteService(favRepo, userRepo, articleRepo)

	user := stubUser(1, "alice", "alice@example.com")
	userRepo.On("FindByID", uint(1)).Return(user, nil)
	articleRepo.On("FindBySlug", "ghost").Return(nil, ErrArticleNotFound)

	err := svc.Unfavorite(1, "ghost")
	assert.ErrorIs(t, err, ErrArticleNotFound)
}
