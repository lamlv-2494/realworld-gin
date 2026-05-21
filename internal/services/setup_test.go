package services

import (
	"os"
	"testing"

	"realworld-gin/internal/models"

	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret")
	os.Exit(m.Run())
}

// ---------------------------------------------------------------------------
// MockUserRepository
// ---------------------------------------------------------------------------

type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) CreateUser(user *models.User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *models.User, data map[string]any) error {
	return m.Called(user, data).Error(0)
}

func (m *MockUserRepository) Follow(follower *models.User, target *models.User) error {
	return m.Called(follower, target).Error(0)
}

func (m *MockUserRepository) Unfollow(follower *models.User, target *models.User) error {
	return m.Called(follower, target).Error(0)
}

func (m *MockUserRepository) IsFollowing(followerID, targetID uint) bool {
	return m.Called(followerID, targetID).Bool(0)
}

// ---------------------------------------------------------------------------
// MockArticleRepository
// ---------------------------------------------------------------------------

type MockArticleRepository struct{ mock.Mock }

func (m *MockArticleRepository) CreateArticle(article *models.Article) error {
	return m.Called(article).Error(0)
}

func (m *MockArticleRepository) FindOrCreateTag(name string) (*models.Tag, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *MockArticleRepository) FindBySlug(slug string) (*models.Article, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Article), args.Error(1)
}

func (m *MockArticleRepository) UpdateArticle(article *models.Article, data map[string]any) error {
	return m.Called(article, data).Error(0)
}

func (m *MockArticleRepository) Delete(article *models.Article) error {
	return m.Called(article).Error(0)
}

func (m *MockArticleRepository) ListArticles(tag, author, favorited string, limit, offset int) ([]*models.Article, int64, error) {
	args := m.Called(tag, author, favorited, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleRepository) FeedArticle(followerID uint, limit, offset int) ([]*models.Article, int64, error) {
	args := m.Called(followerID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Article), args.Get(1).(int64), args.Error(2)
}

func (m *MockArticleRepository) GetTags() ([]string, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// ---------------------------------------------------------------------------
// MockCommentRepository
// ---------------------------------------------------------------------------

type MockCommentRepository struct{ mock.Mock }

func (m *MockCommentRepository) CreateComment(comment *models.Comment) error {
	return m.Called(comment).Error(0)
}

func (m *MockCommentRepository) GetCommentsByArticleID(articleID uint) ([]*models.Comment, error) {
	args := m.Called(articleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Comment), args.Error(1)
}

func (m *MockCommentRepository) GetCommentByID(id uint) (*models.Comment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *MockCommentRepository) DeleteComment(comment *models.Comment) error {
	return m.Called(comment).Error(0)
}

// ---------------------------------------------------------------------------
// MockFavoriteRepository
// ---------------------------------------------------------------------------

type MockFavoriteRepository struct{ mock.Mock }

func (m *MockFavoriteRepository) FavoriteArticle(user *models.User, article *models.Article) error {
	return m.Called(user, article).Error(0)
}

func (m *MockFavoriteRepository) UnfavoriteArticle(user *models.User, article *models.Article) error {
	return m.Called(user, article).Error(0)
}

func (m *MockFavoriteRepository) IsFavorited(userID, articleID uint) bool {
	return m.Called(userID, articleID).Bool(0)
}

func (m *MockFavoriteRepository) GetFavoritesCount(articleID uint) int64 {
	args := m.Called(articleID)
	return args.Get(0).(int64)
}
