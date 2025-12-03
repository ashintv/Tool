package middlewares

import (
	"aetrix/observer/internals/config"
	"aetrix/observer/internals/lib"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type userMiddleware struct {
	cnf *config.Config
}

func NewUserMiddleware(cnf *config.Config) *userMiddleware {
	return &userMiddleware{
		cnf: cnf,
	}
}
func (m *userMiddleware) UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		token, err := lib.ParseToken(tokenString, []byte(m.cnf.USER_JWT_SECRET))
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
		c.Set("user_id", claims["user_id"])
		c.Next()
	}

}
