package controller

import (
	"net/http"
	"strconv"
	"url-shortener/appErrors"
	"url-shortener/middleware"
	"url-shortener/service"

	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsController(analyticsService service.AnalyticsService) *AnalyticsController {
	return &AnalyticsController{
		analyticsService: analyticsService,
	}
}

func (ac *AnalyticsController) GetAnalytics(c *gin.Context) {
	userId, ok := middleware.GetCurrentUserID(c)
	if !ok {
		c.Error(appErrors.ErrUnauthorized)
		return
	}

	urlIdStr := c.Param("id")
	urlId, err := strconv.ParseUint(urlIdStr, 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequestError("Geçersiz URL ID"))
		return
	}

	analytics, err := ac.analyticsService.GetAnalytics(c.Request.Context(), userId, uint(urlId))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, analytics)
}
