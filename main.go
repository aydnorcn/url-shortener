package main

import (
	"log"
	"url-shortener/config"
	"url-shortener/models"
)

func main() {

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

	err = db.AutoMigrate(&models.URL{}, &models.User{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database migrated successfully")
}
