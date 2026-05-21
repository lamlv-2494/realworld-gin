package routes

import (
	"realworld-gin/internal/controllers"
	"realworld-gin/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserGroup(api *gin.RouterGroup, userController *controllers.UserController) {
	userOtional := api.Group("/users")
	{
		userOtional.POST("", userController.Register)
		userOtional.POST("/login", userController.Login)
	}

	userStrict := api.Group("/user", middlewares.AuthMiddleware())
	{
		userStrict.GET("", userController.GetCurrentUser)
		userStrict.PUT("", userController.UpdateCurrentUser)
	}
}

// RegisterProfileRoutes registers profile-related routes (optional auth)
func RegisterProfileRoutes(api *gin.RouterGroup, profileController *controllers.ProfileController) {
	profileOptional := api.Group("/profiles", middlewares.OptionalAuthMiddleware())
	{
		profileOptional.GET("/:username", profileController.GetProfile)
	}

	profilesStrict := api.Group("/profiles", middlewares.AuthMiddleware())
	{
		profilesStrict.POST("/:username/follow", profileController.Follow)
		profilesStrict.DELETE("/:username/follow", profileController.Unfollow)
	}
}

func ResigterArticleRoutes(api *gin.RouterGroup,
	articleController *controllers.ArticleController,
	favoriteController *controllers.FavoriteController,
	commentController *controllers.CommentController,
) {
	articlesOptional := api.Group("/articles", middlewares.OptionalAuthMiddleware())
	{
		articlesOptional.GET("/:slug", articleController.GetBySlug)
		articlesOptional.GET("/:slug/comments", commentController.GetComment)
	}

	articlesStrict := api.Group("/articles", middlewares.AuthMiddleware())
	{
		articlesStrict.POST("", articleController.Create)
		articlesOptional.GET("", articleController.GetArticles)
		articlesStrict.PUT("/:slug", articleController.UpdateArticle)
		articlesStrict.DELETE("/:slug", articleController.DeleteArticle)

		// feed
		articlesStrict.GET("/feed", articleController.GetArticlesFeed)

		// favorite/unfavorite
		articlesStrict.POST("/:slug/favorite", favoriteController.Favorite)
		articlesStrict.DELETE("/:slug/favorite", favoriteController.UnFavorite)

		// Comment
		articlesStrict.POST("/:slug/comments", commentController.CreateComment)
		articlesStrict.DELETE("/:slug/comments/:id", commentController.DeleteComment)
	}
}
