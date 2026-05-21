package services

import (
	"testing"

	"realworld-gin/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCommentService_CreateComment_Success(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewCommentService(commentRepo, articleRepo, userRepo)

	article := stubArticle(10, "my-article", "My Article", 5)
	author := stubUser(2, "bob", "bob@example.com")

	articleRepo.On("FindBySlug", "my-article").Return(article, nil)
	userRepo.On("FindByID", uint(2)).Return(author, nil)
	commentRepo.On("CreateComment", mock.AnythingOfType("*models.Comment")).Return(nil)

	resp, err := svc.CreateComment(2, "my-article", "Great post!")
	require.NoError(t, err)
	assert.Equal(t, "Great post!", resp.Body)
	assert.Equal(t, "bob", resp.Author.Profile.Username)

	commentRepo.AssertExpectations(t)
	articleRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestCommentService_CreateComment_ArticleNotFound(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewCommentService(commentRepo, articleRepo, userRepo)

	articleRepo.On("FindBySlug", "ghost").Return(nil, ErrArticleNotFound)

	_, err := svc.CreateComment(1, "ghost", "body")
	assert.ErrorIs(t, err, ErrArticleNotFound)
}

func TestCommentService_CreateComment_UserNotFound(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewCommentService(commentRepo, articleRepo, userRepo)

	article := stubArticle(10, "article-x", "Article X", 5)
	articleRepo.On("FindBySlug", "article-x").Return(article, nil)
	userRepo.On("FindByID", uint(99)).Return(nil, ErrUserNotFound)

	_, err := svc.CreateComment(99, "article-x", "body")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

// ---------------------------------------------------------------------------
// GetComments
// ---------------------------------------------------------------------------

func TestCommentService_GetComments_Success(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewCommentService(commentRepo, articleRepo, userRepo)

	article := stubArticle(10, "my-article", "Article", 5)
	articleRepo.On("FindBySlug", "my-article").Return(article, nil)

	commenter := models.User{Username: "commenter"}
	commenter.ID = 7
	comments := []*models.Comment{
		{Body: "First", ArticleID: 10, AuthorID: 7, Author: commenter},
		{Body: "Second", ArticleID: 10, AuthorID: 7, Author: commenter},
	}
	commentRepo.On("GetCommentsByArticleID", uint(10)).Return(comments, nil)
	// currentUserID=0 → no IsFollowing call

	resp, err := svc.GetComments(0, "my-article")
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "First", resp[0].Body)

	commentRepo.AssertExpectations(t)
	articleRepo.AssertExpectations(t)
}

func TestCommentService_GetComments_WithFollowing(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewCommentService(commentRepo, articleRepo, userRepo)

	article := stubArticle(10, "my-article", "Article", 5)
	articleRepo.On("FindBySlug", "my-article").Return(article, nil)

	commenter := models.User{Username: "commenter"}
	commenter.ID = 7
	comments := []*models.Comment{{Body: "Hi", ArticleID: 10, AuthorID: 7, Author: commenter}}
	commentRepo.On("GetCommentsByArticleID", uint(10)).Return(comments, nil)
	// currentUserID=2 → IsFollowing is called for each comment author
	userRepo.On("IsFollowing", uint(2), uint(7)).Return(true)

	resp, err := svc.GetComments(2, "my-article")
	require.NoError(t, err)
	require.Len(t, resp, 1)
	assert.True(t, resp[0].Author.Profile.Following)
}

func TestCommentService_GetComments_ArticleNotFound(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	svc := NewCommentService(commentRepo, articleRepo, new(MockUserRepository))

	articleRepo.On("FindBySlug", "ghost").Return(nil, ErrArticleNotFound)

	_, err := svc.GetComments(0, "ghost")
	assert.ErrorIs(t, err, ErrArticleNotFound)
}

// ---------------------------------------------------------------------------
// DeleteComment
// ---------------------------------------------------------------------------

func TestCommentService_DeleteComment_Success(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewCommentService(commentRepo, articleRepo, userRepo)

	user := stubUser(1, "alice", "alice@example.com")
	comment := &models.Comment{Body: "bye", ArticleID: 10, AuthorID: 1}
	comment.ID = 20

	userRepo.On("FindByID", uint(1)).Return(user, nil)
	commentRepo.On("GetCommentByID", uint(20)).Return(comment, nil)
	commentRepo.On("DeleteComment", comment).Return(nil)

	err := svc.DeleteComment(1, "any-slug", 20)
	require.NoError(t, err)
	commentRepo.AssertExpectations(t)
}

func TestCommentService_DeleteComment_UserNotFound(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewCommentService(commentRepo, articleRepo, userRepo)

	userRepo.On("FindByID", uint(99)).Return(nil, ErrUserNotFound)

	err := svc.DeleteComment(99, "slug", 1)
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestCommentService_DeleteComment_CommentNotFound(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewCommentService(commentRepo, articleRepo, userRepo)

	user := stubUser(1, "alice", "alice@example.com")
	userRepo.On("FindByID", uint(1)).Return(user, nil)
	commentRepo.On("GetCommentByID", uint(99)).Return(nil, ErrCommentNotFound)

	err := svc.DeleteComment(1, "slug", 99)
	assert.ErrorIs(t, err, ErrCommentNotFound)
}

func TestCommentService_CreateComment_RepoError(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewCommentService(commentRepo, articleRepo, userRepo)

	article := stubArticle(10, "my-article", "Article", 5)
	author := stubUser(2, "bob", "bob@example.com")
	articleRepo.On("FindBySlug", "my-article").Return(article, nil)
	userRepo.On("FindByID", uint(2)).Return(author, nil)
	commentRepo.On("CreateComment", mock.AnythingOfType("*models.Comment")).Return(ErrCreateComment)

	_, err := svc.CreateComment(2, "my-article", "body")
	assert.ErrorIs(t, err, ErrCreateComment)
}

func TestCommentService_DeleteComment_RepoError(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewCommentService(commentRepo, articleRepo, userRepo)

	user := stubUser(1, "alice", "alice@example.com")
	comment := &models.Comment{Body: "bye", ArticleID: 10, AuthorID: 1}
	comment.ID = 20

	userRepo.On("FindByID", uint(1)).Return(user, nil)
	commentRepo.On("GetCommentByID", uint(20)).Return(comment, nil)
	commentRepo.On("DeleteComment", comment).Return(ErrCreateComment) // any repo error

	err := svc.DeleteComment(1, "slug", 20)
	assert.ErrorIs(t, err, ErrNotCommentAuthor)
}

func TestCommentService_GetComments_RepoError(t *testing.T) {
	commentRepo := new(MockCommentRepository)
	articleRepo := new(MockArticleRepository)
	userRepo := new(MockUserRepository)
	svc := NewCommentService(commentRepo, articleRepo, userRepo)

	article := stubArticle(10, "my-article", "Article", 5)
	articleRepo.On("FindBySlug", "my-article").Return(article, nil)
	commentRepo.On("GetCommentsByArticleID", uint(10)).Return(nil, ErrCommentNotFound)

	_, err := svc.GetComments(0, "my-article")
	assert.ErrorIs(t, err, ErrCommentNotFound)
}
