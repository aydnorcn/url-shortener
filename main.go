package main

import (
	"fmt"
	"log"
	"url-shortener/config"
	"url-shortener/models"
	"url-shortener/routes"
	"url-shortener/validator"
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

	//redisClient := config.NewRedis()

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.URL{}, &models.User{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database migrated successfully")
	router := routes.SetupRouter(db, cfg)

	serverAddr := ":" + cfg.ServerPort
	fmt.Println("Listening on ", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
