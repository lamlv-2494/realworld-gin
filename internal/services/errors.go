package services

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrArticleNotFound     = errors.New("Article not found")
	ErrArticleFeedNotFound = errors.New("Article feed not found")
	ErrCommentNotFound     = errors.New("Comment not found")
	ErrTagsNotFound        = errors.New("Tags not found")

	ErrFollowSelf     = errors.New("Can't follow your self")
	ErrFollowFailed   = errors.New("Can't follow this user")
	ErrUnfollowFailed = errors.New("Can't unfollow this user")

	ErrCreateArticle = errors.New("Failed to create article")
	ErrDeleteArticle = errors.New("Failed to delete article")

	ErrNotArticleAuthor = errors.New("You are not the author of this article")

	ErrCreateComment    = errors.New("Failed to create comment")
	ErrNotCommentAuthor = errors.New("you are not the author of this comment")

	ErrInvalidPassword  = errors.New("Invalid email or password")
	ErrUpdateUser       = errors.New("Failed to update user")
	ErrCreateUserFailed = errors.New("Failed to create user")
)
