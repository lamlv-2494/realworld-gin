package controllers

import (
	"net/http"
	"realworld-gin/internal/services"
	"realworld-gin/internal/utils/constants"

	"github.com/gin-gonic/gin"
)

type FavoriteController struct {
	favoriteService services.FavoriteService
	articleService  services.ArticleService // Tiêm ArticleService để lấy output
}

func NewFavoriteController(f services.FavoriteService, a services.ArticleService) *FavoriteController {
	return &FavoriteController{favoriteService: f, articleService: a}
}

func (ctrl *FavoriteController) Favorite(ctx *gin.Context) {
	currentUserId := ctx.MustGet(constants.UserId).(uint)
	slug := ctx.Param(constants.Slug)

	if err := ctrl.favoriteService.Favorite(currentUserId, slug); err != nil {
		SendError(ctx, err)
		return
	}

	article, err := ctrl.articleService.GetArticle(currentUserId, slug)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, article)
}

func (ctrl *FavoriteController) UnFavorite(ctx *gin.Context) {
	currentUserId := ctx.MustGet(constants.UserId).(uint)
	slug := ctx.Param(constants.Slug)

	if err := ctrl.favoriteService.Unfavorite(currentUserId, slug); err != nil {
		SendError(ctx, err)
		return
	}

	article, err := ctrl.articleService.GetArticle(currentUserId, slug)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, article)
}
