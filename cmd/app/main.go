package main

import (
	"log"
	"os"
	"realworld-gin/internal/config/database"
	"realworld-gin/internal/config/routes"
	"realworld-gin/internal/controllers"
	"realworld-gin/internal/models"
	"realworld-gin/internal/repositories"
	"realworld-gin/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Can't found env file")
	}

	config.ConnectDB()

	err := config.DB.AutoMigrate(&models.User{}, &models.Tag{}, &models.Article{}, &models.Comment{})
	if err != nil {
		log.Fatal("Lỗi Migrate: ", err) // Nếu lỗi nó sẽ dừng app và báo lý do tại đây
	}
	log.Println("✅ Migrate database thành công!")

	userRepo := repositories.NewUserRepository(config.DB)
	artileRepo := repositories.NewArticleRepository(config.DB)
	favoriteRepo := repositories.NewFavoriteRepository(config.DB)
	commentRepo := repositories.NewCommentRepository(config.DB)

	userService := services.NewUserService(userRepo)
	profileService := services.NewProfileService(userRepo)
	articleSerivce := services.NewArticleService(artileRepo, userRepo, favoriteRepo)
	favoriteService := services.NewFavoriteService(favoriteRepo, userRepo, artileRepo)
	commentService := services.NewCommentService(commentRepo, artileRepo, userRepo)

	userController := controllers.NewUserController(userService)
	profileController := controllers.NewProfileController(profileService)
	articleController := controllers.NewArticleController(articleSerivce)
	favoriteController := controllers.NewFavoriteController(favoriteService, articleSerivce)
	commentController := controllers.NewCommentController(commentService)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
		api.GET("/tags", articleController.GetTags)

		routes.RegisterUserGroup(api, userController)
		routes.RegisterProfileRoutes(api, profileController)
		routes.ResigterArticleRoutes(
			api,
			articleController,
			favoriteController,
			commentController,
		)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
