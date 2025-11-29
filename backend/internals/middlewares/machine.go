package middlewares

import (
	"aetrix/observer/internals/lib"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Replace with your actual secret key

func MachineMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		token, err := lib.ParseToken(tokenString , lib.MACHINE_JWT_SECRET)
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("machine_id", claims["machine_id"])
		c.Next()
	}

}
