package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	. "realworld-gin/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SendError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	status := http.StatusInternalServerError
	errMsg := err.Error()

	var validationErrors validator.ValidationErrors
	var jsonSyntaxError *json.SyntaxError
	var jsonTypeError *json.UnmarshalTypeError

	if errors.As(err, &validationErrors) {
		status = http.StatusUnprocessableEntity
		errMsg = ""
		for _, fieldErr := range validationErrors {
			errMsg += fmt.Sprintf("Field '%s' validation error (%s); ", fieldErr.Field(), fieldErr.Tag())
		}

	} else if errors.As(err, &jsonSyntaxError) {
		status = http.StatusBadRequest
		errMsg = "Invalid JSON syntax"
	} else if errors.As(err, &jsonTypeError) {
		status = http.StatusBadRequest
		errMsg = "Invalid JSON type"
	} else {
		switch err {
		case ErrUserNotFound, ErrArticleNotFound, ErrArticleFeedNotFound, ErrCommentNotFound:
			status = http.StatusNotFound
		case ErrNotArticleAuthor, ErrNotCommentAuthor:
			status = http.StatusForbidden
		case ErrFollowSelf, ErrFollowFailed, ErrUnfollowFailed, ErrCreateArticle, ErrDeleteArticle, ErrCreateComment:
			status = http.StatusUnprocessableEntity
		}
	}

	ctx.JSON(status, gin.H{
		"code":  status,
		"error": errMsg,
	})
}
