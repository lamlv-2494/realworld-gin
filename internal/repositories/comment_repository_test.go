package repositories

import (
	"testing"

	"realworld-gin/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentRepository_CreateComment(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	articleRepo := NewArticleRepository(nil)
	repo := NewCommentRepository(nil)

	author := seedUser(t, userRepo, "commenter1", "commenter1@example.com")
	article := seedArticle(t, articleRepo, "comment-article-1", "Comment Article", author.ID)

	comment := &models.Comment{
		Body:      "Great article!",
		ArticleID: article.ID,
		AuthorID:  author.ID,
	}

	err := repo.CreateComment(comment)
	require.NoError(t, err)
	assert.NotZero(t, comment.ID)
}

func TestCommentRepository_GetCommentByID(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	articleRepo := NewArticleRepository(nil)
	repo := NewCommentRepository(nil)

	author := seedUser(t, userRepo, "commenter2", "commenter2@example.com")
	article := seedArticle(t, articleRepo, "comment-article-2", "Comment Article", author.ID)

	comment := &models.Comment{
		Body:      "Nice post",
		ArticleID: article.ID,
		AuthorID:  author.ID,
	}
	require.NoError(t, repo.CreateComment(comment))

	t.Run("found", func(t *testing.T) {
		found, err := repo.GetCommentByID(comment.ID)
		require.NoError(t, err)
		assert.Equal(t, "Nice post", found.Body)
		assert.Equal(t, article.ID, found.ArticleID)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetCommentByID(99999)
		assert.Error(t, err)
	})
}

func TestCommentRepository_GetCommentsByArticleID(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	articleRepo := NewArticleRepository(nil)
	repo := NewCommentRepository(nil)

	author := seedUser(t, userRepo, "commenter3", "commenter3@example.com")
	article := seedArticle(t, articleRepo, "comment-article-3", "Comment Article", author.ID)
	otherArticle := seedArticle(t, articleRepo, "other-article-3", "Other Article", author.ID)

	require.NoError(t, repo.CreateComment(&models.Comment{Body: "First", ArticleID: article.ID, AuthorID: author.ID}))
	require.NoError(t, repo.CreateComment(&models.Comment{Body: "Second", ArticleID: article.ID, AuthorID: author.ID}))
	require.NoError(t, repo.CreateComment(&models.Comment{Body: "Other", ArticleID: otherArticle.ID, AuthorID: author.ID}))

	comments, err := repo.GetCommentsByArticleID(article.ID)
	require.NoError(t, err)
	assert.Len(t, comments, 2, "should only return comments for the requested article")

	for _, c := range comments {
		assert.Equal(t, article.ID, c.ArticleID)
	}
}

func TestCommentRepository_GetCommentsByArticleID_Empty(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	articleRepo := NewArticleRepository(nil)
	repo := NewCommentRepository(nil)

	author := seedUser(t, userRepo, "commenter4", "commenter4@example.com")
	article := seedArticle(t, articleRepo, "comment-article-4", "Empty Comments", author.ID)

	comments, err := repo.GetCommentsByArticleID(article.ID)
	require.NoError(t, err)
	assert.Empty(t, comments)
}

func TestCommentRepository_DeleteComment(t *testing.T) {
	cleanTables(t)
	userRepo := NewUserRepository(nil)
	articleRepo := NewArticleRepository(nil)
	repo := NewCommentRepository(nil)

	author := seedUser(t, userRepo, "commenter5", "commenter5@example.com")
	article := seedArticle(t, articleRepo, "comment-article-5", "Comment Article", author.ID)

	comment := &models.Comment{Body: "Delete me", ArticleID: article.ID, AuthorID: author.ID}
	require.NoError(t, repo.CreateComment(comment))

	err := repo.DeleteComment(comment)
	require.NoError(t, err)

	_, err = repo.GetCommentByID(comment.ID)
	assert.Error(t, err, "comment should not be found after deletion")
}
