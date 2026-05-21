package controllers

import (
	"net/http"
	"testing"

	"realworld-gin/internal/models/dto/responses"
	"realworld-gin/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupArticleRouter(svc services.ArticleService) *gin.Engine {
	r := gin.New()
	ctrl := NewArticleController(svc)
	r.POST("/articles", injectUserID(1), ctrl.Create)
	r.GET("/articles/:slug", ctrl.GetBySlug)
	r.PUT("/articles/:slug", injectUserID(1), ctrl.UpdateArticle)
	r.DELETE("/articles/:slug", injectUserID(1), ctrl.DeleteArticle)
	r.GET("/articles", ctrl.GetArticles)
	r.GET("/articles/feed", injectUserID(1), ctrl.GetArticlesFeed)
	r.GET("/tags", ctrl.GetTags)
	return r
}

func stubArticleResp(slug, title string) *responses.ArticleResponse {
	return &responses.ArticleResponse{Slug: slug, Title: title}
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func TestArticleController_Create_Success(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("CreateArticle", uint(1), "Hello", "desc", "body", []string{"go"}).
		Return(stubArticleResp("hello", "Hello"), nil)

	w := doRequest(r, http.MethodPost, "/articles",
		`{"article":{"title":"Hello","description":"desc","body":"body","tagList":["go"]}}`)

	require.Equal(t, http.StatusCreated, w.Code)
	var resp responses.ArticleResponse
	decodeBody(t, w, &resp)
	assert.Equal(t, "hello", resp.Slug)
	svc.AssertExpectations(t)
}

func TestArticleController_Create_MissingFields(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	// Missing title, description, body → 422
	w := doRequest(r, http.MethodPost, "/articles", `{"article":{}}`)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	svc.AssertNotCalled(t, "CreateArticle")
}

func TestArticleController_Create_ServiceError(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("CreateArticle", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, services.ErrCreateArticle)

	w := doRequest(r, http.MethodPost, "/articles",
		`{"article":{"title":"T","description":"D","body":"B"}}`)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// ---------------------------------------------------------------------------
// GetBySlug
// ---------------------------------------------------------------------------

func TestArticleController_GetBySlug_Success(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("GetArticle", uint(0), "my-slug").
		Return(stubArticleResp("my-slug", "My Article"), nil)

	w := doRequest(r, http.MethodGet, "/articles/my-slug", "")
	require.Equal(t, http.StatusOK, w.Code)
	var resp responses.ArticleResponse
	decodeBody(t, w, &resp)
	assert.Equal(t, "my-slug", resp.Slug)
	svc.AssertExpectations(t)
}

func TestArticleController_GetBySlug_NotFound(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("GetArticle", uint(0), "ghost").Return(nil, services.ErrArticleNotFound)

	w := doRequest(r, http.MethodGet, "/articles/ghost", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ---------------------------------------------------------------------------
// UpdateArticle
// ---------------------------------------------------------------------------

func TestArticleController_UpdateArticle_Success(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("UpdateArticle", uint(1), "my-slug", map[string]any{"body": "updated"}).
		Return(stubArticleResp("my-slug", "Article"), nil)

	w := doRequest(r, http.MethodPut, "/articles/my-slug",
		`{"article":{"body":"updated"}}`)
	require.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestArticleController_UpdateArticle_NotFound(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("UpdateArticle", uint(1), "ghost", mock.Anything).
		Return(nil, services.ErrArticleNotFound)

	w := doRequest(r, http.MethodPut, "/articles/ghost", `{"article":{"body":"x"}}`)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestArticleController_UpdateArticle_NotAuthor(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("UpdateArticle", uint(1), "others-article", mock.Anything).
		Return(nil, services.ErrNotArticleAuthor)

	w := doRequest(r, http.MethodPut, "/articles/others-article", `{"article":{"body":"x"}}`)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ---------------------------------------------------------------------------
// DeleteArticle
// ---------------------------------------------------------------------------

func TestArticleController_DeleteArticle_Success(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("DeleteArticle", uint(1), "bye-article").Return(nil)

	w := doRequest(r, http.MethodDelete, "/articles/bye-article", "")
	assert.Equal(t, http.StatusNoContent, w.Code)
	svc.AssertExpectations(t)
}

func TestArticleController_DeleteArticle_NotFound(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("DeleteArticle", uint(1), "ghost").Return(services.ErrArticleNotFound)

	w := doRequest(r, http.MethodDelete, "/articles/ghost", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestArticleController_DeleteArticle_NotAuthor(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("DeleteArticle", uint(1), "others").Return(services.ErrNotArticleAuthor)

	w := doRequest(r, http.MethodDelete, "/articles/others", "")
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ---------------------------------------------------------------------------
// GetArticles (list)
// ---------------------------------------------------------------------------

func TestArticleController_GetArticles_Success(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	listResp := &responses.ArticleListResponse{
		Articles:   []responses.ArticleResponse{*stubArticleResp("a1", "Article 1")},
		TotalCount: 1,
	}
	svc.On("ListArticles", uint(0), "go", "", "", 5, 1).Return(listResp, nil)

	w := doRequest(r, http.MethodGet, "/articles?tag=go&limit=5&page=1", "")
	require.Equal(t, http.StatusOK, w.Code)
	var resp responses.ArticleListResponse
	decodeBody(t, w, &resp)
	assert.Len(t, resp.Articles, 1)
	svc.AssertExpectations(t)
}

func TestArticleController_GetArticles_InvalidLimitIgnored(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	// non-numeric limit → parsed as 0
	svc.On("ListArticles", uint(0), "", "", "", 0, 0).
		Return(&responses.ArticleListResponse{}, nil)

	w := doRequest(r, http.MethodGet, "/articles?limit=abc", "")
	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// GetArticlesFeed
// ---------------------------------------------------------------------------

func TestArticleController_GetArticlesFeed_Success(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	feedResp := &responses.ArticleListResponse{TotalCount: 2}
	svc.On("FeedArticle", uint(1), 10, 1).Return(feedResp, nil)

	w := doRequest(r, http.MethodGet, "/articles/feed?limit=10&page=1", "")
	require.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestArticleController_GetArticlesFeed_Error(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("FeedArticle", uint(1), 0, 0).Return(nil, services.ErrArticleFeedNotFound)

	w := doRequest(r, http.MethodGet, "/articles/feed", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ---------------------------------------------------------------------------
// GetTags
// ---------------------------------------------------------------------------

func TestArticleController_GetTags_Success(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("GetTags").Return(&responses.TagsResponse{Tags: []string{"go", "rust"}}, nil)

	w := doRequest(r, http.MethodGet, "/tags", "")
	require.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestArticleController_GetTags_Error(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("GetTags").Return(nil, services.ErrTagsNotFound)

	w := doRequest(r, http.MethodGet, "/tags", "")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ---------------------------------------------------------------------------
// Additional tests for missing branches
// ---------------------------------------------------------------------------

// GetBySlug: logged-in user branch (userId = id.(uint))
func TestArticleController_GetBySlug_LoggedIn(t *testing.T) {
	svc := new(MockArticleService)
	r := gin.New()
	ctrl := NewArticleController(svc)
	r.GET("/articles/:slug", injectUserID(1), ctrl.GetBySlug)

	svc.On("GetArticle", uint(1), "my-slug").Return(stubArticleResp("my-slug", "Title"), nil)

	w := doRequest(r, http.MethodGet, "/articles/my-slug", "")
	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

// UpdateArticle: malformed JSON → ShouldBindJSON error → jsonSyntaxError branch in SendError
func TestArticleController_UpdateArticle_InvalidBody(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	w := doRequest(r, http.MethodPut, "/articles/my-slug", `{not valid json}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// UpdateArticle: wrong JSON type → jsonTypeError branch in SendError
func TestArticleController_UpdateArticle_WrongType(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	// title is *string but we send a number → json.UnmarshalTypeError
	w := doRequest(r, http.MethodPut, "/articles/my-slug", `{"article":{"title":123}}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// GetArticles: service error path
func TestArticleController_GetArticles_ServiceError(t *testing.T) {
	svc := new(MockArticleService)
	r := setupArticleRouter(svc)

	svc.On("ListArticles", uint(0), "", "", "", 0, 0).Return(nil, services.ErrArticleNotFound)

	w := doRequest(r, http.MethodGet, "/articles", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// GetArticles: logged-in user branch (userId = id.(uint))
func TestArticleController_GetArticles_LoggedIn(t *testing.T) {
	svc := new(MockArticleService)
	r := gin.New()
	ctrl := NewArticleController(svc)
	r.GET("/articles", injectUserID(1), ctrl.GetArticles)

	svc.On("ListArticles", uint(1), "", "", "", 0, 0).Return(&responses.ArticleListResponse{}, nil)

	w := doRequest(r, http.MethodGet, "/articles", "")
	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}
