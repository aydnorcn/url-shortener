package controller

import (
	"net/http"
	"strconv"
	"url-shortener/appErrors"
	"url-shortener/dto"
	"url-shortener/middleware"
	"url-shortener/service"
	"url-shortener/validator"

	"github.com/gin-gonic/gin"
)

type UrlController struct {
	urlService service.UrlService
}

func NewUrlController(urlService service.UrlService) *UrlController {
	return &UrlController{urlService: urlService}
}

func (u *UrlController) CreateUrl(c *gin.Context) {
	var req dto.CreateUrlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.ErrInvalidJSON)
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		c.Error(err)
		return
	}

	userId, ok := middleware.GetCurrentUserID(c)
	if !ok {
		c.Error(appErrors.ErrUnauthorized)
		return
	}

	url, err := u.urlService.CreateUrl(userId, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, url)
}

func (u *UrlController) GetUrl(c *gin.Context) {
	userId, ok := middleware.GetCurrentUserID(c)
	if !ok {
		c.Error(appErrors.ErrUnauthorized)
		return
	}

	urlId := c.Param("id")
	urlIdInt, err := strconv.ParseUint(urlId, 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequestError("Geçersiz URL ID"))
		return
	}

	url, err := u.urlService.GetUrl(userId, uint(urlIdInt))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, url)
}

func (u *UrlController) GetAllUrls(c *gin.Context) {
	userId, ok := middleware.GetCurrentUserID(c)
	if !ok {
		c.Error(appErrors.ErrUnauthorized)
		return
	}

	urls, err := u.urlService.GetUserUrls(userId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, urls)
}

func (u *UrlController) UpdateUrl(c *gin.Context) {
	var req dto.UpdateUrlRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.ErrInvalidJSON)
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		c.Error(err)
		return
	}

	userId, ok := middleware.GetCurrentUserID(c)
	if !ok {
		c.Error(appErrors.ErrUnauthorized)
		return
	}

	urlId := c.Param("id")
	id, err := strconv.ParseUint(urlId, 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequestError("Geçersiz URL ID"))
		return
	}

	url, err := u.urlService.UpdateUrl(userId, uint(id), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, url)
}

func (u *UrlController) DeleteUrl(c *gin.Context) {
	userId, ok := middleware.GetCurrentUserID(c)
	if !ok {
		c.Error(appErrors.ErrUnauthorized)
		return
	}

	urlId := c.Param("id")
	id, err := strconv.ParseUint(urlId, 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequestError("Geçersiz URL ID"))
		return
	}

	if err := u.urlService.DeleteUrl(userId, uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (u *UrlController) ActivateUrl(c *gin.Context) {
	userId, ok := middleware.GetCurrentUserID(c)
	if !ok {
		c.Error(appErrors.ErrUnauthorized)
		return
	}

	urlId := c.Param("id")
	id, err := strconv.ParseUint(urlId, 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequestError("Geçersiz URL ID"))
		return
	}

	if err := u.urlService.ActivateUrl(userId, uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (u *UrlController) DeactivateUrl(c *gin.Context) {
	userId, ok := middleware.GetCurrentUserID(c)
	if !ok {
		c.Error(appErrors.ErrUnauthorized)
		return
	}

	urlId := c.Param("id")
	id, err := strconv.ParseUint(urlId, 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequestError("Geçersiz URL ID"))
		return
	}

	if err := u.urlService.DeactivateUrl(userId, uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (u *UrlController) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")

	originalUrl, err := u.urlService.Redirect(shortCode)
	if err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusFound, originalUrl)
}
