package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFavoriteRepository_FavoriteAndUnfavorite(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	articleRepo := NewArticleRepository(nil)
	repo := NewFavoriteRepository(nil)

	user := seedUser(t, userRepo, "favuser1", "favuser1@example.com")
	article := seedArticle(t, articleRepo, "fav-article-1", "Fav Article", user.ID)

	t.Run("favorite article", func(t *testing.T) {
		err := repo.FavoriteArticle(user, article)
		require.NoError(t, err)
		assert.True(t, repo.IsFavorited(user.ID, article.ID))
	})

	t.Run("favorite count is 1", func(t *testing.T) {
		count := repo.GetFavoritesCount(article.ID)
		assert.EqualValues(t, 1, count)
	})

	t.Run("unfavorite article", func(t *testing.T) {
		err := repo.UnfavoriteArticle(user, article)
		require.NoError(t, err)
		assert.False(t, repo.IsFavorited(user.ID, article.ID))
	})

	t.Run("favorite count is 0 after unfavorite", func(t *testing.T) {
		count := repo.GetFavoritesCount(article.ID)
		assert.EqualValues(t, 0, count)
	})
}

func TestFavoriteRepository_IsFavorited_MultipleUsers(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	articleRepo := NewArticleRepository(nil)
	repo := NewFavoriteRepository(nil)

	author := seedUser(t, userRepo, "favauthor2", "favauthor2@example.com")
	user1 := seedUser(t, userRepo, "favuser2a", "favuser2a@example.com")
	user2 := seedUser(t, userRepo, "favuser2b", "favuser2b@example.com")
	article := seedArticle(t, articleRepo, "fav-article-2", "Multi Fav", author.ID)

	require.NoError(t, repo.FavoriteArticle(user1, article))

	assert.True(t, repo.IsFavorited(user1.ID, article.ID))
	assert.False(t, repo.IsFavorited(user2.ID, article.ID), "user2 has not favorited")
}

func TestFavoriteRepository_GetFavoritesCount_MultipleUsers(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	articleRepo := NewArticleRepository(nil)
	repo := NewFavoriteRepository(nil)

	author := seedUser(t, userRepo, "favauthor3", "favauthor3@example.com")
	user1 := seedUser(t, userRepo, "favuser3a", "favuser3a@example.com")
	user2 := seedUser(t, userRepo, "favuser3b", "favuser3b@example.com")
	article := seedArticle(t, articleRepo, "fav-article-3", "Count Article", author.ID)

	require.NoError(t, repo.FavoriteArticle(user1, article))
	require.NoError(t, repo.FavoriteArticle(user2, article))

	count := repo.GetFavoritesCount(article.ID)
	assert.EqualValues(t, 2, count)
}

func TestFavoriteRepository_IsFavorited_NotFavorited(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	articleRepo := NewArticleRepository(nil)
	repo := NewFavoriteRepository(nil)

	user := seedUser(t, userRepo, "favuser4", "favuser4@example.com")
	article := seedArticle(t, articleRepo, "fav-article-4", "No Fav", user.ID)

	assert.False(t, repo.IsFavorited(user.ID, article.ID))
	assert.EqualValues(t, 0, repo.GetFavoritesCount(article.ID))
}
