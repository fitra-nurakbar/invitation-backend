package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Hanya izinkan role tertentu
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		role := userRole.(string)
		for _, allowed := range roles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Akses ditolak — role tidak diizinkan",
		})
	}
}