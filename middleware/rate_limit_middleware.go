package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Counter struct {
	Count     int
	StartTime time.Time
}

type RateLimiter struct {
	MaxRequestPerWindow int
	RateLimitWindow     time.Duration
	mu                  sync.Mutex
	ipRequest           map[string]*Counter
}

func NewRateLimiter(requests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		MaxRequestPerWindow: requests,
		RateLimitWindow:     window,
		ipRequest:           make(map[string]*Counter),
	}
}

func (r *RateLimiter) LimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		r.mu.Lock()

		info, exists := r.ipRequest[ip]
		if !exists {
			info = &Counter{Count: 1, StartTime: now}
			r.ipRequest[ip] = info
		} else {
			if now.Sub(info.StartTime) < r.RateLimitWindow {
				info.Count++
			} else {
				info.Count = 1
				info.StartTime = now
			}
		}
		r.mu.Unlock()

		if info.Count > r.MaxRequestPerWindow {
			c.AbortWithStatusJSON(
				429, gin.H{
					"code":    429,
					"message": "Too many requests",
				})
		}
		c.Next()
	}
}
