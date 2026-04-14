package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var Secret = []byte("super-secret")

func CreateToken(user string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
	})

	return token.SignedString(Secret)
}

func Validate(tokenStr string) (string, bool) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return Secret, nil
	})

	if err != nil || !token.Valid {
		return "", false
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["user"].(string), true
}
