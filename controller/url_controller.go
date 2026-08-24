package controller

import (
	"math"
	"net/http"
	"strconv"
	"time"
	"url-shortener/appErrors"
	"url-shortener/dto"
	"url-shortener/middleware"
	"url-shortener/service"
	"url-shortener/utils"
	"url-shortener/validator"
	"url-shortener/worker"

	"github.com/gin-gonic/gin"
)

type UrlController struct {
	urlService service.UrlService
	worker     worker.AnalyticsWorker
}

func NewUrlController(urlService service.UrlService, worker worker.AnalyticsWorker) *UrlController {
	return &UrlController{
		urlService: urlService,
		worker:     worker,
	}
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	userId, ok := middleware.GetCurrentUserID(c)
	if !ok {
		c.Error(appErrors.ErrUnauthorized)
		return
	}

	urls, total, err := u.urlService.GetUserUrls(userId, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": urls,
		"pagination": gin.H{
			"page":        page,
			"size":        pageSize,
			"total":       total,
			"total_pages": math.Ceil(float64(total) / float64(pageSize)),
		},
	})
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

	url, err := u.urlService.Redirect(shortCode, c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	// Record click event asynchronously through worker pool
	if u.worker != nil {
		event := dto.ClickEvent{
			URLID:     url.ID,
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Referer:   c.Request.Referer(),
			Country:   utils.ParseCountry(c),
			Device:    utils.ParseDevice(c.Request.UserAgent()),
			Timestamp: time.Now(),
		}
		u.worker.Process(event)
	}

	c.Redirect(http.StatusFound, url.OriginalURL)
}
