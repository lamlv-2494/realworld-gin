package services

import (
	"testing"

	"realworld-gin/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// stubArticle returns a fully populated Article that service layer can read fields from.
func stubArticle(id uint, slug, title string, authorID uint) *models.Article {
	a := &models.Article{
		Slug:        slug,
		Title:       title,
		Description: "desc",
		Body:        "body",
		AuthorID:    authorID,
		Author:      models.User{Username: "author", Bio: "", Image: ""},
		Tags:        []*models.Tag{{Name: "go"}},
	}
	a.ID = id
	return a
}

func stubUser(id uint, username, email string) *models.User {
	u := &models.User{Username: username, Email: email}
	u.ID = id
	return u
}

// ---------------------------------------------------------------------------
// CreateArticle
// ---------------------------------------------------------------------------

func TestArticleService_CreateArticle_Success(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	favRepo := new(MockFavoriteRepository)
	svc := NewArticleService(articleRepo, userRepo, favRepo)

	author := stubUser(1, "alice", "alice@example.com")
	userRepo.On("FindByID", uint(1)).Return(author, nil)
	articleRepo.On("FindOrCreateTag", "go").Return(&models.Tag{Name: "go"}, nil)
	articleRepo.On("CreateArticle", mock.AnythingOfType("*models.Article")).Return(nil)

	resp, err := svc.CreateArticle(1, "Hello World", "desc", "body", []string{"go"})
	require.NoError(t, err)
	assert.Equal(t, "Hello World", resp.Title)
	assert.Equal(t, "alice", resp.Author.Profile.Username)
	assert.ElementsMatch(t, []string{"go"}, resp.TagList)

	articleRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestArticleService_CreateArticle_AuthorNotFound(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	favRepo := new(MockFavoriteRepository)
	svc := NewArticleService(articleRepo, userRepo, favRepo)

	userRepo.On("FindByID", uint(99)).Return(nil, ErrUserNotFound)

	_, err := svc.CreateArticle(99, "Title", "desc", "body", nil)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

// ---------------------------------------------------------------------------
// GetArticle
// ---------------------------------------------------------------------------

func TestArticleService_GetArticle_Success_Anonymous(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	favRepo := new(MockFavoriteRepository)
	svc := NewArticleService(articleRepo, userRepo, favRepo)

	article := stubArticle(10, "test-slug", "Test", 1)
	articleRepo.On("FindBySlug", "test-slug").Return(article, nil)
	// currentUserID=0 → IsFavorited is NOT called by the service
	favRepo.On("GetFavoritesCount", uint(10)).Return(int64(3))

	resp, err := svc.GetArticle(0, "test-slug")
	require.NoError(t, err)
	assert.Equal(t, "test-slug", resp.Slug)
	assert.False(t, resp.Favorited)
	assert.EqualValues(t, 3, resp.FavoritesCount)

	articleRepo.AssertExpectations(t)
	favRepo.AssertExpectations(t)
}

func TestArticleService_GetArticle_WithFavoriteAndFollow(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	favRepo := new(MockFavoriteRepository)
	svc := NewArticleService(articleRepo, userRepo, favRepo)

	article := stubArticle(10, "test-slug", "Test", 5)
	articleRepo.On("FindBySlug", "test-slug").Return(article, nil)
	userRepo.On("IsFollowing", uint(2), uint(5)).Return(true)
	favRepo.On("IsFavorited", uint(2), uint(10)).Return(true)
	favRepo.On("GetFavoritesCount", uint(10)).Return(int64(7))

	resp, err := svc.GetArticle(2, "test-slug")
	require.NoError(t, err)
	assert.True(t, resp.Favorited)
	assert.True(t, resp.Author.Profile.Following)
	assert.EqualValues(t, 7, resp.FavoritesCount)
}

func TestArticleService_GetArticle_NotFound(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	articleRepo.On("FindBySlug", "no-such-slug").Return(nil, ErrArticleNotFound)

	_, err := svc.GetArticle(0, "no-such-slug")
	assert.ErrorIs(t, err, ErrArticleNotFound)
}

// ---------------------------------------------------------------------------
// UpdateArticle
// ---------------------------------------------------------------------------

func TestArticleService_UpdateArticle_Success(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	favRepo := new(MockFavoriteRepository)
	svc := NewArticleService(articleRepo, userRepo, favRepo)

	article := stubArticle(10, "old-slug", "Old Title", 1)
	updated := stubArticle(10, "new-title-xxxx", "New Title", 1)

	articleRepo.On("FindBySlug", "old-slug").Return(article, nil)
	articleRepo.On("UpdateArticle", article, mock.Anything).Return(nil)
	// The new slug has a random suffix from slug_gen, so match any string != "old-slug"
	articleRepo.On("FindBySlug", mock.MatchedBy(func(s string) bool { return s != "old-slug" })).Return(updated, nil)
	userRepo.On("IsFollowing", uint(1), uint(1)).Return(false)
	favRepo.On("IsFavorited", uint(1), uint(10)).Return(false)
	favRepo.On("GetFavoritesCount", uint(10)).Return(int64(0))

	resp, err := svc.UpdateArticle(1, "old-slug", map[string]any{"title": "New Title"})
	require.NoError(t, err)
	assert.Equal(t, "New Title", resp.Title)
	articleRepo.AssertExpectations(t)
}

func TestArticleService_UpdateArticle_NotFound(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	articleRepo.On("FindBySlug", "ghost").Return(nil, ErrArticleNotFound)

	_, err := svc.UpdateArticle(1, "ghost", map[string]any{"body": "x"})
	assert.ErrorIs(t, err, ErrArticleNotFound)
}

func TestArticleService_UpdateArticle_NotAuthor(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	article := stubArticle(10, "my-article", "Title", 5) // authorID=5
	articleRepo.On("FindBySlug", "my-article").Return(article, nil)

	_, err := svc.UpdateArticle(99, "my-article", map[string]any{"body": "x"}) // userID=99
	assert.ErrorIs(t, err, ErrNotArticleAuthor)
}

// ---------------------------------------------------------------------------
// DeleteArticle
// ---------------------------------------------------------------------------

func TestArticleService_DeleteArticle_Success(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	article := stubArticle(10, "bye-article", "Bye", 1)
	articleRepo.On("FindBySlug", "bye-article").Return(article, nil)
	articleRepo.On("Delete", article).Return(nil)

	err := svc.DeleteArticle(1, "bye-article")
	require.NoError(t, err)
	articleRepo.AssertExpectations(t)
}

func TestArticleService_DeleteArticle_NotFound(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	articleRepo.On("FindBySlug", "ghost").Return(nil, ErrArticleNotFound)

	err := svc.DeleteArticle(1, "ghost")
	assert.ErrorIs(t, err, ErrArticleNotFound)
}

func TestArticleService_DeleteArticle_NotAuthor(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	article := stubArticle(10, "others-article", "Title", 5) // authorID=5
	articleRepo.On("FindBySlug", "others-article").Return(article, nil)

	err := svc.DeleteArticle(99, "others-article") // userID=99
	assert.ErrorIs(t, err, ErrNotArticleAuthor)
}

// ---------------------------------------------------------------------------
// ListArticles
// ---------------------------------------------------------------------------

func TestArticleService_ListArticles_Success(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	favRepo := new(MockFavoriteRepository)
	svc := NewArticleService(articleRepo, userRepo, favRepo)

	articles := []*models.Article{
		stubArticle(1, "a1", "Article 1", 10),
		stubArticle(2, "a2", "Article 2", 10),
	}
	userRepo.On("IsFollowing", uint(0), uint(10)).Return(false).Maybe()
	articleRepo.On("ListArticles", "", "", "", 10, 0).Return(articles, int64(2), nil)

	resp, err := svc.ListArticles(0, "", "", "", 10, 1)
	require.NoError(t, err)
	assert.EqualValues(t, 2, resp.TotalCount)
	assert.Len(t, resp.Articles, 2)
	articleRepo.AssertExpectations(t)
}

func TestArticleService_ListArticles_DefaultPagination(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	// page=0 and limit=0 should default to page=1, limit=10 → offset=0
	articleRepo.On("ListArticles", "", "", "", 10, 0).Return([]*models.Article{}, int64(0), nil)

	resp, err := svc.ListArticles(0, "", "", "", 0, 0)
	require.NoError(t, err)
	assert.Empty(t, resp.Articles)
	articleRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// FeedArticle
// ---------------------------------------------------------------------------

func TestArticleService_FeedArticle_Success(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	articles := []*models.Article{stubArticle(1, "feed-1", "Feed", 5)}
	// page=1, limit=5 → offset=0
	articleRepo.On("FeedArticle", uint(3), 5, 0).Return(articles, int64(1), nil)

	resp, err := svc.FeedArticle(3, 5, 1)
	require.NoError(t, err)
	assert.EqualValues(t, 1, resp.TotalCount)
	assert.Len(t, resp.Articles, 1)
	articleRepo.AssertExpectations(t)
}

func TestArticleService_FeedArticle_Error(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	articleRepo.On("FeedArticle", uint(3), 10, 0).Return(nil, int64(0), ErrArticleFeedNotFound)

	_, err := svc.FeedArticle(3, 10, 1)
	assert.ErrorIs(t, err, ErrArticleFeedNotFound)
}

// ---------------------------------------------------------------------------
// GetTags
// ---------------------------------------------------------------------------

func TestArticleService_GetTags_Success(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	articleRepo.On("GetTags").Return([]string{"go", "rust", "python"}, nil)

	resp, err := svc.GetTags()
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"go", "rust", "python"}, resp.Tags)
	articleRepo.AssertExpectations(t)
}

func TestArticleService_GetTags_Error(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	articleRepo.On("GetTags").Return(nil, ErrTagsNotFound)

	_, err := svc.GetTags()
	assert.ErrorIs(t, err, ErrTagsNotFound)
}

// ---------------------------------------------------------------------------
// Additional tests to reach 100% coverage
// ---------------------------------------------------------------------------

func TestArticleService_CreateArticle_RepoError(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewArticleService(articleRepo, userRepo, new(MockFavoriteRepository))

	author := stubUser(1, "alice", "alice@example.com")
	userRepo.On("FindByID", uint(1)).Return(author, nil)
	articleRepo.On("CreateArticle", mock.AnythingOfType("*models.Article")).Return(ErrCreateArticle)

	_, err := svc.CreateArticle(1, "Hello", "desc", "body", nil)
	assert.ErrorIs(t, err, ErrCreateArticle)
}

func TestArticleService_UpdateArticle_RepoError(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	article := stubArticle(10, "my-slug", "Title", 1)
	articleRepo.On("FindBySlug", "my-slug").Return(article, nil)
	articleRepo.On("UpdateArticle", article, mock.Anything).Return(ErrCreateArticle) // any repo error

	_, err := svc.UpdateArticle(1, "my-slug", map[string]any{"body": "x"})
	assert.ErrorIs(t, err, ErrNotArticleAuthor)
}

func TestArticleService_DeleteArticle_RepoError(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	article := stubArticle(10, "my-slug", "Title", 1)
	articleRepo.On("FindBySlug", "my-slug").Return(article, nil)
	articleRepo.On("Delete", article).Return(ErrDeleteArticle)

	err := svc.DeleteArticle(1, "my-slug")
	assert.ErrorIs(t, err, ErrDeleteArticle)
}

func TestArticleService_ListArticles_Error(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	articleRepo.On("ListArticles", "", "", "", 10, 0).Return(nil, int64(0), ErrArticleNotFound)

	_, err := svc.ListArticles(0, "", "", "", 10, 1)
	assert.ErrorIs(t, err, ErrArticleNotFound)
}

func TestArticleService_ListArticles_WithLoggedInUser(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewArticleService(articleRepo, userRepo, new(MockFavoriteRepository))

	articles := []*models.Article{stubArticle(1, "a1", "Article", 5)}
	articleRepo.On("ListArticles", "", "", "", 10, 0).Return(articles, int64(1), nil)
	userRepo.On("IsFollowing", uint(2), uint(5)).Return(true)

	resp, err := svc.ListArticles(2, "", "", "", 10, 1)
	require.NoError(t, err)
	require.Len(t, resp.Articles, 1)
	assert.True(t, resp.Articles[0].Author.Profile.Following)
}

func TestArticleService_FeedArticle_DefaultPagination(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	// page=0, limit=0 → defaults: page=1 (offset=0), limit=10
	articleRepo.On("FeedArticle", uint(1), 10, 0).Return([]*models.Article{}, int64(0), nil)

	resp, err := svc.FeedArticle(1, 0, 0)
	require.NoError(t, err)
	assert.Empty(t, resp.Articles)
	articleRepo.AssertExpectations(t)
}

func TestArticleService_FeedArticle_EmptyResult(t *testing.T) {
	articleRepo := new(MockArticleRepository)
	svc := NewArticleService(articleRepo, new(MockUserRepository), new(MockFavoriteRepository))

	// Empty articles slice (no error) → articlesResponse initialised as empty slice
	articleRepo.On("FeedArticle", uint(1), 5, 0).Return([]*models.Article{}, int64(0), nil)

	resp, err := svc.FeedArticle(1, 5, 1)
	require.NoError(t, err)
	assert.Empty(t, resp.Articles)
	assert.EqualValues(t, 0, resp.TotalCount)
}
