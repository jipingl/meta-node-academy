package middleware

import (
	"errors"
	"net/http"

	"example.com/blog/internal/config"
	"example.com/blog/internal/util"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusOK, config.Fail(http.StatusUnauthorized, errors.New("Unauthorized")))
			c.Abort()
			return
		}
		claims, err := util.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusOK, config.Fail(http.StatusUnauthorized, errors.New("Unauthorized")))
			c.Abort()
			return
		}
		c.Set(config.USER_ID, claims.UserID)
		c.Next()
	}
}
