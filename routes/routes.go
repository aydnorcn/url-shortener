package routes

import (
	"time"
	"url-shortener/cache"
	"url-shortener/config"
	"url-shortener/controller"
	"url-shortener/middleware"
	"url-shortener/repository"
	"url-shortener/service"
	"url-shortener/worker"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg config.Config, cache cache.Cache, analyticsWorker worker.AnalyticsWorker) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.MetricsMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(middleware.NewRateLimiter(60, time.Minute).LimitMiddleware())

	// Serve frontend website
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	userRepo := repository.NewUserRepository(db)
	urlRepo := repository.NewUrlRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	authService := service.NewAuthService(userRepo, &cfg)
	urlService := service.NewUrlService(urlRepo, cache)
	analyticsService := service.NewAnalyticsService(analyticsRepo, urlRepo)

	authController := controller.NewAuthController(authService)
	urlController := controller.NewUrlController(urlService, analyticsWorker)
	analyticsController := controller.NewAnalyticsController(analyticsService)

	// Root redirect endpoint
	r.GET("/:shortCode", urlController.Redirect)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "UP",
				"service": "URL SHORTENER SERVICE",
				"version": "1.0.0",
			})
		})

		// V1 API routes
		v1 := api.Group("/v1")
		{
			authV1 := v1.Group("/auth")
			{
				authV1.POST("/login", authController.Login)
				authV1.POST("/register", authController.Register)
			}

			urlV1 := v1.Group("/urls")
			urlV1.Use(middleware.AuthMiddleware(cfg.JwtSecret))
			{
				urlV1.POST("", urlController.CreateUrl)
				urlV1.POST("/", urlController.CreateUrl)
				urlV1.GET("", urlController.GetAllUrls)
				urlV1.GET("/", urlController.GetAllUrls)
				urlV1.GET("/:id", urlController.GetUrl)
				urlV1.PATCH("/:id", urlController.UpdateUrl)
				urlV1.DELETE("/:id", urlController.DeleteUrl)
				urlV1.PATCH("/:id/activate", urlController.ActivateUrl)
				urlV1.PATCH("/:id/deactivate", urlController.DeactivateUrl)
				urlV1.GET("/:id/analytics", analyticsController.GetAnalytics)
			}
		}

		// Backward compatibility for /api/...
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authController.Login)
			authGroup.POST("/register", authController.Register)
		}
		urlGroup := api.Group("/urls")
		urlGroup.Use(middleware.AuthMiddleware(cfg.JwtSecret))
		{
			urlGroup.POST("", urlController.CreateUrl)
			urlGroup.POST("/", urlController.CreateUrl)
			urlGroup.GET("", urlController.GetAllUrls)
			urlGroup.GET("/", urlController.GetAllUrls)
			urlGroup.GET("/:id", urlController.GetUrl)
			urlGroup.PATCH("/:id", urlController.UpdateUrl)
			urlGroup.DELETE("/:id", urlController.DeleteUrl)
			urlGroup.PATCH("/:id/activate", urlController.ActivateUrl)
			urlGroup.PATCH("/:id/deactivate", urlController.DeactivateUrl)
			urlGroup.GET("/:id/analytics", analyticsController.GetAnalytics)
		}

		api.GET("/:shortCode", urlController.Redirect)
	}

	return r
}
