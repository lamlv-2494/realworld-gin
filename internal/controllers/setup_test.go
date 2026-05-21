package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"realworld-gin/internal/models/dto/requests"
	"realworld-gin/internal/models/dto/responses"
	"realworld-gin/internal/utils/constants"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

// injectUserID is a middleware stub that puts a fixed userID into the context.
func injectUserID(id uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constants.UserId, id)
		c.Next()
	}
}

// decodeBody unmarshals the recorder body into v.
func decodeBody(t *testing.T, w *httptest.ResponseRecorder, v any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), v))
}

// ---------------------------------------------------------------------------
// MockUserService
// ---------------------------------------------------------------------------

type MockUserService struct{ mock.Mock }

func (m *MockUserService) Register(req requests.RegisterRequest) (*responses.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.UserResponse), args.Error(1)
}

func (m *MockUserService) Login(req requests.LoginRequest) (*responses.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.UserResponse), args.Error(1)
}

func (m *MockUserService) GetCurrentUser(id uint) (*responses.UserResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.UserResponse), args.Error(1)
}

func (m *MockUserService) UpdateCurrentUser(id uint, data map[string]interface{}) (*responses.UserResponse, error) {
	args := m.Called(id, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.UserResponse), args.Error(1)
}

// ---------------------------------------------------------------------------
// MockArticleService
// ---------------------------------------------------------------------------

type MockArticleService struct{ mock.Mock }

func (m *MockArticleService) CreateArticle(authorID uint, title, description, body string, tagList []string) (*responses.ArticleResponse, error) {
	args := m.Called(authorID, title, description, body, tagList)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.ArticleResponse), args.Error(1)
}

func (m *MockArticleService) GetArticle(currentUserID uint, slug string) (*responses.ArticleResponse, error) {
	args := m.Called(currentUserID, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.ArticleResponse), args.Error(1)
}

func (m *MockArticleService) UpdateArticle(currentUserID uint, slug string, data map[string]any) (*responses.ArticleResponse, error) {
	args := m.Called(currentUserID, slug, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.ArticleResponse), args.Error(1)
}

func (m *MockArticleService) DeleteArticle(currentUserID uint, slug string) error {
	return m.Called(currentUserID, slug).Error(0)
}

func (m *MockArticleService) ListArticles(currentUserID uint, tag, author, favorited string, limit, page int) (*responses.ArticleListResponse, error) {
	args := m.Called(currentUserID, tag, author, favorited, limit, page)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.ArticleListResponse), args.Error(1)
}

func (m *MockArticleService) FeedArticle(currentUserID uint, limit, page int) (*responses.ArticleListResponse, error) {
	args := m.Called(currentUserID, limit, page)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.ArticleListResponse), args.Error(1)
}

func (m *MockArticleService) GetTags() (*responses.TagsResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.TagsResponse), args.Error(1)
}

// ---------------------------------------------------------------------------
// MockCommentService
// ---------------------------------------------------------------------------

type MockCommentService struct{ mock.Mock }

func (m *MockCommentService) CreateComment(currentUserID uint, slug, body string) (*responses.CommentResponse, error) {
	args := m.Called(currentUserID, slug, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.CommentResponse), args.Error(1)
}

func (m *MockCommentService) GetComments(currentUserID uint, slug string) ([]responses.CommentResponse, error) {
	args := m.Called(currentUserID, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]responses.CommentResponse), args.Error(1)
}

func (m *MockCommentService) DeleteComment(currentUserID uint, slug string, commentID uint) error {
	return m.Called(currentUserID, slug, commentID).Error(0)
}

// ---------------------------------------------------------------------------
// MockProfileService
// ---------------------------------------------------------------------------

type MockProfileService struct{ mock.Mock }

func (m *MockProfileService) GetProfile(currentUserID uint, targetUsername string) (*responses.ProfileData, error) {
	args := m.Called(currentUserID, targetUsername)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.ProfileData), args.Error(1)
}

func (m *MockProfileService) Follow(currentUserID uint, targetUsername string) (*responses.ProfileResponse, error) {
	args := m.Called(currentUserID, targetUsername)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.ProfileResponse), args.Error(1)
}

func (m *MockProfileService) Unfollow(currentUserID uint, targetUsername string) (*responses.ProfileResponse, error) {
	args := m.Called(currentUserID, targetUsername)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.ProfileResponse), args.Error(1)
}

// ---------------------------------------------------------------------------
// MockFavoriteService
// ---------------------------------------------------------------------------

type MockFavoriteService struct{ mock.Mock }

func (m *MockFavoriteService) Favorite(currentUserID uint, slug string) error {
	return m.Called(currentUserID, slug).Error(0)
}

func (m *MockFavoriteService) Unfavorite(currentUserID uint, slug string) error {
	return m.Called(currentUserID, slug).Error(0)
}
