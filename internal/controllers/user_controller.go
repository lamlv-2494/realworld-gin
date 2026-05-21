package controllers

import (
	"net/http"
	"realworld-gin/internal/models/dto/requests"
	"realworld-gin/internal/services"
	"realworld-gin/internal/utils/constants"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService}
}

func (ctrl *UserController) Register(ctx *gin.Context) {
	var req requests.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {

		SendError(ctx, err)
		return
	}

	user, err := ctrl.userService.Register(req)
	if err != nil {
		SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (ctrl *UserController) Login(c *gin.Context) {
	var req requests.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, err)
		return
	}

	user, err := ctrl.userService.Login(req)
	if err != nil {
		SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) GetCurrentUser(c *gin.Context) {
	userID := c.MustGet(constants.UserId).(uint)
	token := strings.Split(c.GetHeader(constants.Authorization), " ")[1]

	user, err := ctrl.userService.GetCurrentUser(userID)
	if err != nil {
		SendError(c, err)
		return
	}

	user.User.Token = token

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) UpdateCurrentUser(c *gin.Context) {
	userID := c.MustGet(constants.UserId).(uint)

	var req requests.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, err)
		return
	}

	updateData := make(map[string]any)
	if req.User.Username != nil {
		updateData["username"] = *req.User.Username
	}
	if req.User.Email != nil {
		updateData["email"] = *req.User.Email
	}
	if req.User.Password != nil {
		updateData["password"] = *req.User.Password
	}
	if req.User.Bio != nil {
		updateData["bio"] = *req.User.Bio
	}
	if req.User.Image != nil {
		updateData["image"] = *req.User.Image
	}

	response, err := ctrl.userService.UpdateCurrentUser(userID, updateData)
	if err != nil {
		SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
