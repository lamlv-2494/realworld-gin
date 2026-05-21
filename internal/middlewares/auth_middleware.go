package middlewares

import (
	"net/http"

	"realworld-gin/internal/utils"
	"realworld-gin/internal/utils/constants"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(constants.Authorization)
		if authHeader == "" {
			abort(ctx, "Request Authorization!")
			return
		}

		part := strings.Split(authHeader, " ")
		if len(part) != 2 || part[0] != constants.Token {
			abort(ctx, "Invalid token!")
			return
		}

		userID, err := utils.ValidateToken(part[1])
		if err != nil {
			abort(ctx, err.Error())
			return
		}

		ctx.Set(constants.UserId, userID)
		ctx.Next()

	}
}

func abort(ctx *gin.Context, errorMessage any) {
	errorMessages := gin.H{
		"code":  http.StatusUnauthorized,
		"error": errorMessage,
	}
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorMessages)
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(constants.Authorization)
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")

			if len(parts) == 2 && parts[0] == constants.Token {
				userID, err := utils.ValidateToken(parts[1])
				if err == nil {
					ctx.Set(constants.UserId, userID)
				}
			}
		}

		ctx.Next()
	}
}
