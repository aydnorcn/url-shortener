package routes

import (
	"url-shortener/controller"
	"url-shortener/repository"
	"url-shortener/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	//TODO: User and url repo
	userRepo := repository.NewUserRepository(db)
	urlRepo := repository.NewUrlRepository(db)

	//TODO: User and url service
	authService := service.NewAuthService(userRepo)
	urlService := service.NewUrlService(urlRepo)

	//TODO: User and url controller
	authController := controller.NewAuthController(authService)
	urlController := controller.NewUrlController(urlService)

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "UP",
				"service": "URL SHORTENER SERVICE",
				"version": "1.0.0",
			})
		})

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authController.Login)
			authGroup.POST("/register", authController.Register)
		}
		urlGroup := api.Group("/urls")
		{
			urlGroup.POST("/", urlController.CreateUrl)
			urlGroup.GET("/:id", urlController.GetUrl)
			urlGroup.PUT("/:id", urlController.UpdateUrl)
			urlGroup.DELETE("/:id", urlController.DeleteUrl)
		}

		api.GET("/:shortCode", urlController.GetUrl)
	}

	return r
}
