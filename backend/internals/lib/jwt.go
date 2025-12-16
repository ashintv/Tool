package lib

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, JWT_SECRET []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECRET)
}

func ParseToken(tokenString string, JWT_SECRET []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenInvalidClaims
		}
		return JWT_SECRET, nil
	})
}
