package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	corsAllowMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	corsAllowHeaders = "Authorization, Content-Type, Accept, Origin, X-Requested-With"
	corsMaxAge       = "43200"
)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowed[origin] = struct{}{}
		}
	}
	allowAll := len(allowed) == 0

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		allowedOrigin := ""
		if allowAll {
			allowedOrigin = "*"
		} else if _, ok := allowed[origin]; ok {
			allowedOrigin = origin
		}

		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			if allowedOrigin != "*" {
				c.Header("Vary", "Origin")
			}
			c.Header("Access-Control-Expose-Headers", "X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset, Retry-After")
		}

		if c.Request.Method == http.MethodOptions {
			if allowedOrigin == "" {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			requestedHeaders := strings.TrimSpace(c.GetHeader("Access-Control-Request-Headers"))
			if requestedHeaders == "" {
				requestedHeaders = corsAllowHeaders
			}

			c.Header("Access-Control-Allow-Methods", corsAllowMethods)
			c.Header("Access-Control-Allow-Headers", requestedHeaders)
			c.Header("Access-Control-Max-Age", corsMaxAge)
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
