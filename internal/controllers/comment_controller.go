package controllers

import (
	"net/http"

	"realworld-gin/internal/models/dto/requests"
	"realworld-gin/internal/services"
	"realworld-gin/internal/utils/constants"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	service services.CommentService
}

func NewCommentController(service services.CommentService) *CommentController {
	return &CommentController{service: service}
}

func (ctrl *CommentController) CreateComment(ctx *gin.Context) {
	currentUserID := ctx.MustGet(constants.UserId).(uint)

	slug := ctx.Param(constants.Slug)
	var req requests.CommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendError(ctx, err)
		return
	}

	response, err := ctrl.service.CreateComment(currentUserID, slug, req.Comment.Body)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response)

}

func (ctrl *CommentController) GetComment(ctx *gin.Context) {
	var currentUserID uint = 0
	if id, exists := ctx.Get(constants.UserId); exists {
		currentUserID = id.(uint)
	}

	slug := ctx.Param(constants.Slug)

	response, err := ctrl.service.GetComments(currentUserID, slug)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (ctrl *CommentController) DeleteComment(ctx *gin.Context) {
	currentUserID := ctx.MustGet(constants.UserId).(uint)
	slug := ctx.Param(constants.Slug)
	id, err := strconv.ParseUint(ctx.Param(constants.Id), 10, 64)
	if err != nil {
		SendError(ctx, err)
		return
	}

	if err := ctrl.service.DeleteComment(currentUserID, slug, uint(id)); err != nil {
		SendError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
