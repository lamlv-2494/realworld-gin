package controllers

import (
	"net/http"
	"realworld-gin/internal/models/dto/responses"
	"realworld-gin/internal/services"
	"realworld-gin/internal/utils/constants"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	profileService services.ProfileService
}

func NewProfileController(profileService services.ProfileService) *ProfileController {
	return &ProfileController{profileService}
}

func (ctrl *ProfileController) GetProfile(ctx *gin.Context) {
	userName := ctx.Param(constants.Username)

	var currentUserID uint = 0
	if id, exists := ctx.Get(constants.UserId); exists {
		currentUserID = id.(uint)
	}

	profile, err := ctrl.profileService.GetProfile(currentUserID, userName)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, responses.ProfileResponse{Profile: *profile})

}

func (ctrl *ProfileController) Follow(ctx *gin.Context) {
	userName := ctx.Param(constants.Username)
	currentUserID := ctx.MustGet(constants.UserId).(uint)

	profile, err := ctrl.profileService.Follow(currentUserID, userName)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func (ctrl *ProfileController) Unfollow(ctx *gin.Context) {
	userName := ctx.Param(constants.Username)
	currentUserID := ctx.MustGet(constants.UserId).(uint)

	profile, err := ctrl.profileService.Unfollow(currentUserID, userName)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, profile)

}
