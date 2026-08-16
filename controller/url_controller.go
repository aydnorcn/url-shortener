package controller

import (
	"url-shortener/service"

	"github.com/gin-gonic/gin"
)

type UrlController struct {
	urlService service.UrlService
}

func NewUrlController(urlService service.UrlService) *UrlController {
	return &UrlController{urlService: urlService}
}

func (u *UrlController) CreateUrl(c *gin.Context) {

}

func (u *UrlController) GetUrl(c *gin.Context) {

}

func (u *UrlController) UpdateUrl(c *gin.Context) {

}

func (u *UrlController) DeleteUrl(c *gin.Context) {

}

func (u *UrlController) Redirect(c *gin.Context) {

}
