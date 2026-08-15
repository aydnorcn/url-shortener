package routes

import (
	"url-shortener/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	//TODO: User and url repo
	userRepo := repository.NewUserRepository(db)
	urlRepo := repository.NewUrlRepository(db)
	//TODO: User and url service

	//TODO: User and url controller

	return r
}
