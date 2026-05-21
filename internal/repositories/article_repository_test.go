package repositories

import (
	"testing"

	"realworld-gin/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// seedArticle creates a minimal article owned by the given authorID.
func seedArticle(t *testing.T, repo ArticleRepository, slug, title string, authorID uint, tags ...*models.Tag) *models.Article {
	t.Helper()
	article := &models.Article{
		Slug:        slug,
		Title:       title,
		Description: "desc",
		Body:        "body",
		AuthorID:    authorID,
		Tags:        tags,
	}
	require.NoError(t, repo.CreateArticle(article))
	return article
}

func TestArticleRepository_CreateArticle(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	repo := NewArticleRepository(nil)

	author := seedUser(t, userRepo, "author1", "author1@example.com")

	article := &models.Article{
		Slug:     "create-test",
		Title:    "Create Test",
		Body:     "body",
		AuthorID: author.ID,
	}
	err := repo.CreateArticle(article)
	require.NoError(t, err)
	assert.NotZero(t, article.ID)
}

func TestArticleRepository_FindOrCreateTag(t *testing.T) {
	cleanTables(t)
	repo := NewArticleRepository(nil)

	t.Run("creates new tag", func(t *testing.T) {
		tag, err := repo.FindOrCreateTag("golang")
		require.NoError(t, err)
		assert.NotZero(t, tag.ID)
		assert.Equal(t, "golang", tag.Name)
	})

	t.Run("finds existing tag without duplicate", func(t *testing.T) {
		tag1, err := repo.FindOrCreateTag("golang")
		require.NoError(t, err)

		tag2, err := repo.FindOrCreateTag("golang")
		require.NoError(t, err)

		assert.Equal(t, tag1.ID, tag2.ID, "same tag should be returned on second call")
	})
}

func TestArticleRepository_FindBySlug(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	repo := NewArticleRepository(nil)

	author := seedUser(t, userRepo, "author2", "author2@example.com")
	seedArticle(t, repo, "find-by-slug", "Find By Slug", author.ID)

	t.Run("found", func(t *testing.T) {
		found, err := repo.FindBySlug("find-by-slug")
		require.NoError(t, err)
		assert.Equal(t, "Find By Slug", found.Title)
		assert.Equal(t, author.ID, found.AuthorID)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindBySlug("non-existent-slug")
		assert.Error(t, err)
	})
}

func TestArticleRepository_UpdateArticle(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	repo := NewArticleRepository(nil)

	author := seedUser(t, userRepo, "author3", "author3@example.com")
	article := seedArticle(t, repo, "update-test", "Original Title", author.ID)

	err := repo.UpdateArticle(article, map[string]any{
		"title":       "Updated Title",
		"description": "Updated description",
	})
	require.NoError(t, err)

	updated, err := repo.FindBySlug("update-test")
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, "Updated description", updated.Description)
}

func TestArticleRepository_Delete(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	repo := NewArticleRepository(nil)

	author := seedUser(t, userRepo, "author4", "author4@example.com")
	article := seedArticle(t, repo, "delete-test", "Delete Me", author.ID)

	err := repo.Delete(article)
	require.NoError(t, err)

	_, err = repo.FindBySlug("delete-test")
	assert.Error(t, err, "article should not be found after deletion")
}

func TestArticleRepository_ListArticles_NoFilter(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	repo := NewArticleRepository(nil)

	author := seedUser(t, userRepo, "author5", "author5@example.com")
	seedArticle(t, repo, "list-1", "Article One", author.ID)
	seedArticle(t, repo, "list-2", "Article Two", author.ID)

	articles, count, err := repo.ListArticles("", "", "", 10, 0)
	require.NoError(t, err)
	assert.EqualValues(t, 2, count)
	assert.Len(t, articles, 2)
}

func TestArticleRepository_ListArticles_FilterByAuthor(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	repo := NewArticleRepository(nil)

	author := seedUser(t, userRepo, "author6", "author6@example.com")
	other := seedUser(t, userRepo, "other6", "other6@example.com")
	seedArticle(t, repo, "a6-1", "Author Article", author.ID)
	seedArticle(t, repo, "o6-1", "Other Article", other.ID)

	articles, count, err := repo.ListArticles("", "author6", "", 10, 0)
	require.NoError(t, err)
	assert.EqualValues(t, 1, count)
	assert.Len(t, articles, 1)
	assert.Equal(t, "Author Article", articles[0].Title)
}

func TestArticleRepository_ListArticles_FilterByTag(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	repo := NewArticleRepository(nil)

	author := seedUser(t, userRepo, "author7", "author7@example.com")

	goTag, err := repo.FindOrCreateTag("go")
	require.NoError(t, err)

	seedArticle(t, repo, "tagged-article", "Tagged Article", author.ID, goTag)
	seedArticle(t, repo, "untagged-article", "Untagged Article", author.ID)

	articles, count, err := repo.ListArticles("go", "", "", 10, 0)
	require.NoError(t, err)
	assert.EqualValues(t, 1, count)
	assert.Len(t, articles, 1)
	assert.Equal(t, "Tagged Article", articles[0].Title)
}

func TestArticleRepository_ListArticles_Pagination(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	repo := NewArticleRepository(nil)

	author := seedUser(t, userRepo, "author8", "author8@example.com")
	for i := range 5 {
		seedArticle(t, repo, "page-article-"+string(rune('a'+i)), "Paginated", author.ID)
	}

	articles, count, err := repo.ListArticles("", "", "", 2, 0)
	require.NoError(t, err)
	assert.EqualValues(t, 5, count)
	assert.Len(t, articles, 2, "limit should be respected")

	articles2, _, err := repo.ListArticles("", "", "", 2, 2)
	require.NoError(t, err)
	assert.Len(t, articles2, 2, "offset should skip first 2")
}

func TestArticleRepository_FeedArticle(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	repo := NewArticleRepository(nil)

	follower := seedUser(t, userRepo, "follower9", "follower9@example.com")
	followed := seedUser(t, userRepo, "followed9", "followed9@example.com")
	unrelated := seedUser(t, userRepo, "unrelated9", "unrelated9@example.com")

	require.NoError(t, userRepo.Follow(follower, followed))

	seedArticle(t, repo, "feed-article-1", "Feed Article", followed.ID)
	seedArticle(t, repo, "unrelated-article", "Unrelated", unrelated.ID)

	articles, count, err := repo.FeedArticle(follower.ID, 10, 0)
	require.NoError(t, err)
	assert.EqualValues(t, 1, count)
	assert.Len(t, articles, 1)
	assert.Equal(t, "Feed Article", articles[0].Title)
}

func TestArticleRepository_GetTags(t *testing.T) {
	cleanTables(t)
	repo := NewArticleRepository(nil)

	_, err := repo.FindOrCreateTag("tagA")
	require.NoError(t, err)
	_, err = repo.FindOrCreateTag("tagB")
	require.NoError(t, err)

	tags, err := repo.GetTags()
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"tagA", "tagB"}, tags)
}
