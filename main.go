package main

import (
	"context"
	"fmt"
	"log"
	"url-shortener/cache"
	"url-shortener/config"
	"url-shortener/metrics"
	"url-shortener/models"
	"url-shortener/repository"
	"url-shortener/routes"
	"url-shortener/service"
	"url-shortener/validator"
	"url-shortener/worker"
)

func main() {
	// Initialize validation engine
	validator.Init()

	cfg := config.Load()

	db, err := config.Connect(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	if err != nil {
		log.Fatal(err)
	}

	redisClient := config.NewRedis()
	redisCache := cache.NewRedisCache(redisClient)
	metrics.Init()

	// Auto-migrate models including URLClick
	err = db.AutoMigrate(&models.URL{}, &models.User{}, &models.URLClick{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database migrated successfully")

	// Initialize repositories and services needed for analytics worker
	urlRepo := repository.NewUrlRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	analyticsService := service.NewAnalyticsService(analyticsRepo, urlRepo)

	// Initialize and start analytics worker pool
	analyticsWorker := worker.NewAnalyticsWorker(analyticsService, 5, 1000)
	analyticsWorker.Start(context.Background())
	defer analyticsWorker.Stop()

	router := routes.SetupRouter(db, cfg, redisCache, analyticsWorker)

	serverAddr := ":" + cfg.ServerPort
	fmt.Println("Listening on ", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
