package controllers

import (
	"net/http"
	. "realworld-gin/internal/models/dto/requests"
	"realworld-gin/internal/services"
	"realworld-gin/internal/utils/constants"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	articleService services.ArticleService
}

func NewArticleController(service services.ArticleService) *ArticleController {
	return &ArticleController{articleService: service}
}

func (ctrl *ArticleController) Create(ctx *gin.Context) {
	userId := ctx.MustGet(constants.UserId).(uint)

	var req CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendError(ctx, err)
		return
	}

	article, err := ctrl.articleService.CreateArticle(userId, req.Article.Title, req.Article.Description, req.Article.Body, req.Article.TagList)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, article)
}

func (ctrl *ArticleController) GetBySlug(ctx *gin.Context) {
	var userId uint = 0
	if id, exist := ctx.Get(constants.UserId); exist {
		userId = id.(uint)
	}

	slug := ctx.Param(constants.Slug)

	article, err := ctrl.articleService.GetArticle(userId, slug)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, article)
}

func (ctrl *ArticleController) UpdateArticle(ctx *gin.Context) {
	userId := ctx.MustGet(constants.UserId).(uint)

	slug := ctx.Param(constants.Slug)

	var req UpdateArticleRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendError(ctx, err)
		return
	}

	updateData := make(map[string]any)
	if req.Article.Title != nil {
		updateData[constants.Title] = *req.Article.Title
	}
	if req.Article.Description != nil {
		updateData[constants.Description] = *req.Article.Description
	}
	if req.Article.Body != nil {
		updateData[constants.Body] = *req.Article.Body
	}

	updatedArticle, err := ctrl.articleService.UpdateArticle(userId, slug, updateData)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, updatedArticle)
}

func (ctrl *ArticleController) DeleteArticle(ctx *gin.Context) {
	userId := ctx.MustGet(constants.UserId).(uint)
	slug := ctx.Param(constants.Slug)

	if err := ctrl.articleService.DeleteArticle(userId, slug); err != nil {
		SendError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (ctrl *ArticleController) GetArticles(ctx *gin.Context) {
	var userId uint = 0
	if id, exists := ctx.Get(constants.UserId); exists {
		userId = id.(uint)
	}

	tag := ctx.Query(constants.Tag)
	author := ctx.Query(constants.Author)
	favorited := ctx.Query(constants.Favorite)
	limit := 0
	page := 0
	if parsedLimit, err := strconv.Atoi(ctx.Query(constants.Limit)); err == nil {
		limit = parsedLimit
	}
	if parsedPage, err := strconv.Atoi(ctx.Query(constants.Page)); err == nil {
		page = parsedPage
	}

	response, err := ctrl.articleService.ListArticles(userId, tag, author, favorited, limit, page)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (ctrl *ArticleController) GetArticlesFeed(ctx *gin.Context) {
	userId := ctx.MustGet(constants.UserId).(uint)

	limit := 0
	page := 0
	if parsedLimit, err := strconv.Atoi(ctx.Query(constants.Limit)); err == nil {
		limit = parsedLimit
	}
	if parsedPage, err := strconv.Atoi(ctx.Query(constants.Page)); err == nil {
		page = parsedPage
	}

	response, err := ctrl.articleService.FeedArticle(userId, limit, page)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (ctrl *ArticleController) GetTags(ctx *gin.Context) {
	tags, err := ctrl.articleService.GetTags()
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, tags)
}
