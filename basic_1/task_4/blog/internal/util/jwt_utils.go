package util

import (
	"time"

	"example.com/blog/internal/entity"
	"github.com/dgrijalva/jwt-go"
)

var secretKey = "temporary-secret-key"

type CustomClaims struct {
	UserID   uint   `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(user *entity.User) (string, int64, error) {
	expAt := time.Now().Add(time.Hour * 24).Unix()
	claims := CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", 0, err
	}
	return tokenStr, expAt, nil
}

func ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
