package controllers

import (
	"net/http"
	"testing"

	"realworld-gin/internal/models/dto/responses"
	"realworld-gin/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupFavoriteRouter(favSvc services.FavoriteService, articleSvc services.ArticleService) *gin.Engine {
	r := gin.New()
	ctrl := NewFavoriteController(favSvc, articleSvc)
	r.POST("/articles/:slug/favorite", injectUserID(1), ctrl.Favorite)
	r.DELETE("/articles/:slug/favorite", injectUserID(1), ctrl.UnFavorite)
	return r
}

// ---------------------------------------------------------------------------
// Favorite
// ---------------------------------------------------------------------------

func TestFavoriteController_Favorite_Success(t *testing.T) {
	favSvc := new(MockFavoriteService)
	articleSvc := new(MockArticleService)
	r := setupFavoriteRouter(favSvc, articleSvc)

	favSvc.On("Favorite", uint(1), "cool-article").Return(nil)
	articleSvc.On("GetArticle", uint(1), "cool-article").
		Return(&responses.ArticleResponse{Slug: "cool-article", Favorited: true, FavoritesCount: 1}, nil)

	w := doRequest(r, http.MethodPost, "/articles/cool-article/favorite", "")
	require.Equal(t, http.StatusOK, w.Code)

	var resp responses.ArticleResponse
	decodeBody(t, w, &resp)
	assert.True(t, resp.Favorited)
	assert.EqualValues(t, 1, resp.FavoritesCount)

	favSvc.AssertExpectations(t)
	articleSvc.AssertExpectations(t)
}

func TestFavoriteController_Favorite_ArticleNotFound(t *testing.T) {
	favSvc := new(MockFavoriteService)
	articleSvc := new(MockArticleService)
	r := setupFavoriteRouter(favSvc, articleSvc)

	favSvc.On("Favorite", uint(1), "ghost").Return(services.ErrArticleNotFound)

	w := doRequest(r, http.MethodPost, "/articles/ghost/favorite", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
	articleSvc.AssertNotCalled(t, "GetArticle")
}

func TestFavoriteController_Favorite_GetArticleError(t *testing.T) {
	favSvc := new(MockFavoriteService)
	articleSvc := new(MockArticleService)
	r := setupFavoriteRouter(favSvc, articleSvc)

	favSvc.On("Favorite", uint(1), "some-article").Return(nil)
	articleSvc.On("GetArticle", uint(1), "some-article").Return(nil, services.ErrArticleNotFound)

	w := doRequest(r, http.MethodPost, "/articles/some-article/favorite", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ---------------------------------------------------------------------------
// UnFavorite
// ---------------------------------------------------------------------------

func TestFavoriteController_UnFavorite_Success(t *testing.T) {
	favSvc := new(MockFavoriteService)
	articleSvc := new(MockArticleService)
	r := setupFavoriteRouter(favSvc, articleSvc)

	favSvc.On("Unfavorite", uint(1), "cool-article").Return(nil)
	articleSvc.On("GetArticle", uint(1), "cool-article").
		Return(&responses.ArticleResponse{Slug: "cool-article", Favorited: false, FavoritesCount: 0}, nil)

	w := doRequest(r, http.MethodDelete, "/articles/cool-article/favorite", "")
	require.Equal(t, http.StatusOK, w.Code)

	var resp responses.ArticleResponse
	decodeBody(t, w, &resp)
	assert.False(t, resp.Favorited)
	assert.EqualValues(t, 0, resp.FavoritesCount)

	favSvc.AssertExpectations(t)
	articleSvc.AssertExpectations(t)
}

func TestFavoriteController_UnFavorite_UserNotFound(t *testing.T) {
	favSvc := new(MockFavoriteService)
	articleSvc := new(MockArticleService)
	r := setupFavoriteRouter(favSvc, articleSvc)

	favSvc.On("Unfavorite", uint(1), "some-article").Return(services.ErrUserNotFound)

	w := doRequest(r, http.MethodDelete, "/articles/some-article/favorite", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
	articleSvc.AssertNotCalled(t, "GetArticle")
}

// UnFavorite: Unfavorite succeeds but GetArticle fails → second error branch
func TestFavoriteController_UnFavorite_GetArticleError(t *testing.T) {
	favSvc := new(MockFavoriteService)
	articleSvc := new(MockArticleService)
	r := setupFavoriteRouter(favSvc, articleSvc)

	favSvc.On("Unfavorite", uint(1), "some-article").Return(nil)
	articleSvc.On("GetArticle", uint(1), "some-article").Return(nil, services.ErrArticleNotFound)

	w := doRequest(r, http.MethodDelete, "/articles/some-article/favorite", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
	favSvc.AssertExpectations(t)
	articleSvc.AssertExpectations(t)
}
