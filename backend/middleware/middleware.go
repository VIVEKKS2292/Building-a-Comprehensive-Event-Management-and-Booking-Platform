package middleware

import (
	"net/http"

	"backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
			c.Abort()
			return
		}

		// Parse the token and get claims
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user ID from claims
		userID, ok := claims["userID"].(float64) // JWT claims typically store numbers as float64
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token structure"})
			c.Abort()
			return
		}

		// Check role if roles are specified
		if len(roles) > 0 {
			userRole, ok := claims["role"].(string)
			if !ok {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token structure"})
				c.Abort()
				return
			}

			roleAllowed := false
			for _, role := range roles {
				if userRole == role {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort()
				return
			}
		}

		// Set user ID and claims in the context
		c.Set("userID", uint(userID)) // Convert float64 to uint
		c.Set("claims", claims)

		c.Next()
	}
}
