package middleware

import (
	"net/http"
	"strings"

	"konsera-backend/internal/helpers"
	"konsera-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

const (
	ClaimsKey = "auth_claims"
	UserIDKey = "user_id"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			helpers.Error(c, http.StatusUnauthorized, "Authorization token is required", nil)
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			helpers.Error(c, http.StatusUnauthorized, "Invalid or expired authorization token", nil)
			c.Abort()
			return
		}

		c.Set(ClaimsKey, claims)
		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		value, exists := c.Get(ClaimsKey)
		claims, valid := value.(*utils.JWTClaims)
		if !exists || !valid || claims == nil {
			helpers.Error(c, http.StatusUnauthorized, "Authentication is required", nil)
			c.Abort()
			return
		}

		for _, role := range claims.Roles {
			if _, ok := allowed[role]; ok {
				c.Next()
				return
			}
		}

		helpers.Error(c, http.StatusForbidden, "You do not have permission to access this resource", nil)
		c.Abort()
	}
}
