package middleware

import (
	"net/http"
	"sync"
	"time"

	"konsera-backend/internal/helpers"

	"github.com/gin-gonic/gin"
)

type rateLimitClient struct {
	windowStart time.Time
	requests    int
}

type RateLimiter struct {
	limit   int
	window  time.Duration
	mu      sync.Mutex
	clients map[string]rateLimitClient
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	if limit < 1 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &RateLimiter{limit: limit, window: window, clients: make(map[string]rateLimitClient)}
}

func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		key := c.ClientIP()
		if userID, exists := c.Get(UserIDKey); exists {
			if value, ok := userID.(string); ok && value != "" {
				key = "user:" + value
			}
		}

		r.mu.Lock()
		if len(r.clients) > 10000 {
			for clientKey, value := range r.clients {
				if now.Sub(value.windowStart) >= r.window {
					delete(r.clients, clientKey)
				}
			}
		}
		client, exists := r.clients[key]
		if !exists || now.Sub(client.windowStart) >= r.window {
			client = rateLimitClient{windowStart: now}
		}
		client.requests++
		remaining := r.limit - client.requests
		if remaining < 0 {
			remaining = 0
		}
		resetAt := client.windowStart.Add(r.window)
		r.clients[key] = client
		r.mu.Unlock()

		c.Header("X-RateLimit-Limit", formatInt(r.limit))
		c.Header("X-RateLimit-Remaining", formatInt(remaining))
		c.Header("X-RateLimit-Reset", formatInt(int(resetAt.Unix())))

		if client.requests > r.limit {
			retryAfter := int(time.Until(resetAt).Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", formatInt(retryAfter))
			helpers.Error(c, http.StatusTooManyRequests, "Too many requests", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func formatInt(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	buffer := [20]byte{}
	index := len(buffer)
	for value > 0 {
		index--
		buffer[index] = byte('0' + value%10)
		value /= 10
	}
	if negative {
		index--
		buffer[index] = '-'
	}
	return string(buffer[index:])
}
