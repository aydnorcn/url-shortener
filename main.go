package main

import (
	"log"
	"url-shortener/config"
	"url-shortener/models"
	"url-shortener/routes"
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
	router := routes.SetupRouter(db)

	err = router.Run()
	if err != nil {
		log.Fatal("Error starting server", err)
	}

}
