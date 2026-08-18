package middleware

import (
	"errors"
	"net/http"
	"url-shortener/appErrors"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *appErrors.AppError

		if errors.As(err, &appErr) {
			c.JSON(appErr.Status, appErr)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Internal Server error",
		})
	}
}
