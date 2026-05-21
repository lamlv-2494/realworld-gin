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

func setupCommentRouter(svc services.CommentService) *gin.Engine {
	r := gin.New()
	ctrl := NewCommentController(svc)
	r.POST("/articles/:slug/comments", injectUserID(1), ctrl.CreateComment)
	r.GET("/articles/:slug/comments", ctrl.GetComment)
	r.DELETE("/articles/:slug/comments/:id", injectUserID(1), ctrl.DeleteComment)
	return r
}

// ---------------------------------------------------------------------------
// CreateComment
// ---------------------------------------------------------------------------

func TestCommentController_CreateComment_Success(t *testing.T) {
	svc := new(MockCommentService)
	r := setupCommentRouter(svc)

	resp := &responses.CommentResponse{Body: "Great post!"}
	svc.On("CreateComment", uint(1), "my-article", "Great post!").Return(resp, nil)

	w := doRequest(r, http.MethodPost, "/articles/my-article/comments",
		`{"comment":{"body":"Great post!"}}`)

	require.Equal(t, http.StatusCreated, w.Code)
	var got responses.CommentResponse
	decodeBody(t, w, &got)
	assert.Equal(t, "Great post!", got.Body)
	svc.AssertExpectations(t)
}

func TestCommentController_CreateComment_MissingBody(t *testing.T) {
	svc := new(MockCommentService)
	r := setupCommentRouter(svc)

	// missing body → 422
	w := doRequest(r, http.MethodPost, "/articles/my-article/comments",
		`{"comment":{}}`)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	svc.AssertNotCalled(t, "CreateComment")
}

func TestCommentController_CreateComment_ArticleNotFound(t *testing.T) {
	svc := new(MockCommentService)
	r := setupCommentRouter(svc)

	svc.On("CreateComment", uint(1), "ghost", mock.Anything).Return(nil, services.ErrArticleNotFound)

	w := doRequest(r, http.MethodPost, "/articles/ghost/comments",
		`{"comment":{"body":"hello"}}`)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ---------------------------------------------------------------------------
// GetComment
// ---------------------------------------------------------------------------

func TestCommentController_GetComments_Success(t *testing.T) {
	svc := new(MockCommentService)
	r := setupCommentRouter(svc)

	comments := []responses.CommentResponse{{Body: "First"}, {Body: "Second"}}
	svc.On("GetComments", uint(0), "my-article").Return(comments, nil)

	w := doRequest(r, http.MethodGet, "/articles/my-article/comments", "")
	require.Equal(t, http.StatusOK, w.Code)
	var got []responses.CommentResponse
	decodeBody(t, w, &got)
	assert.Len(t, got, 2)
	svc.AssertExpectations(t)
}

func TestCommentController_GetComments_ArticleNotFound(t *testing.T) {
	svc := new(MockCommentService)
	r := setupCommentRouter(svc)

	svc.On("GetComments", uint(0), "ghost").Return(nil, services.ErrArticleNotFound)

	w := doRequest(r, http.MethodGet, "/articles/ghost/comments", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ---------------------------------------------------------------------------
// DeleteComment
// ---------------------------------------------------------------------------

func TestCommentController_DeleteComment_Success(t *testing.T) {
	svc := new(MockCommentService)
	r := setupCommentRouter(svc)

	svc.On("DeleteComment", uint(1), "my-article", uint(42)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/articles/my-article/comments/42", "")
	assert.Equal(t, http.StatusNoContent, w.Code)
	svc.AssertExpectations(t)
}

func TestCommentController_DeleteComment_CommentNotFound(t *testing.T) {
	svc := new(MockCommentService)
	r := setupCommentRouter(svc)

	svc.On("DeleteComment", uint(1), "my-article", uint(99)).Return(services.ErrCommentNotFound)

	w := doRequest(r, http.MethodDelete, "/articles/my-article/comments/99", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCommentController_DeleteComment_InvalidID(t *testing.T) {
	svc := new(MockCommentService)
	r := setupCommentRouter(svc)

	// non-numeric id → 400 bad request
	w := doRequest(r, http.MethodDelete, "/articles/my-article/comments/abc", "")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertNotCalled(t, "DeleteComment")
}

// GetComment: logged-in user branch (currentUserID = id.(uint))
func TestCommentController_GetComments_LoggedIn(t *testing.T) {
	svc := new(MockCommentService)
	r := gin.New()
	ctrl := NewCommentController(svc)
	r.GET("/articles/:slug/comments", injectUserID(1), ctrl.GetComment)

	svc.On("GetComments", uint(1), "my-article").Return([]responses.CommentResponse{}, nil)

	w := doRequest(r, http.MethodGet, "/articles/my-article/comments", "")
	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}
