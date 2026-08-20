package controller

import (
	"net/http"
	"strconv"
	"url-shortener/dto"
	"url-shortener/middleware"
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
	var req dto.CreateUrlRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := middleware.GetCurrentUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	url, err := u.urlService.CreateUrl(userId, req)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, &url)
}

func (u *UrlController) GetUrl(c *gin.Context) {
	userId, ok := middleware.GetCurrentUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	urlId := c.Param("id")
	urlIdInt, _ := strconv.ParseUint(urlId, 0, 64)
	url, err := u.urlService.GetUrl(userId, uint(urlIdInt))

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, &url)
}

func (u *UrlController) GetAllUrls(c *gin.Context) {
	userId, ok := middleware.GetCurrentUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	urls, err := u.urlService.GetUserUrls(userId)

	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, &urls)
}

func (u *UrlController) UpdateUrl(c *gin.Context) {
	var req dto.UpdateUrlRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := middleware.GetCurrentUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	urlId := c.Param("id")

	asd, _ := strconv.ParseUint(urlId, 0, 64)

	url, err := u.urlService.UpdateUrl(userId, uint(asd), req)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, &url)
}

func (u *UrlController) DeleteUrl(c *gin.Context) {

	urlId := c.Param("id")
	id, err := strconv.ParseUint(urlId, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userId, ok := middleware.GetCurrentUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := u.urlService.DeleteUrl(userId, uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (u *UrlController) ActivateUrl(c *gin.Context) {
	urlId := c.Param("id")
	id, err := strconv.ParseUint(urlId, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userId, ok := middleware.GetCurrentUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := u.urlService.ActivateUrl(userId, uint(id)); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (u *UrlController) DeactivateUrl(c *gin.Context) {
	urlId := c.Param("id")
	id, err := strconv.ParseUint(urlId, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userId, ok := middleware.GetCurrentUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
