package routes

import (
	"time"
	"url-shortener/config"
	"url-shortener/controller"
	"url-shortener/middleware"
	"url-shortener/repository"
	"url-shortener/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg config.Config, redis *redis.Client) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(middleware.NewRateLimiter(6, time.Minute).LimitMiddleware())

	userRepo := repository.NewUserRepository(db)
	urlRepo := repository.NewUrlRepository(db)

	authService := service.NewAuthService(userRepo, &cfg)
	urlService := service.NewUrlService(urlRepo, redis)

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
		urlGroup.Use(middleware.AuthMiddleware(cfg.JwtSecret))
		{
			urlGroup.POST("/", urlController.CreateUrl)
			urlGroup.GET("/", urlController.GetAllUrls)
			urlGroup.GET("/:id", urlController.GetUrl)
			urlGroup.PATCH("/:id", urlController.UpdateUrl)
			urlGroup.DELETE("/:id", urlController.DeleteUrl)
			urlGroup.PATCH("/:id/activate", urlController.ActivateUrl)
			urlGroup.PATCH("/:id/deactivate", urlController.DeactivateUrl)
		}

		api.GET("/:shortCode", urlController.Redirect)
	}

	return r
}
