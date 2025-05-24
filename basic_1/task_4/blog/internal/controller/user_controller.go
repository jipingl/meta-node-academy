package controller

import (
	"net/http"

	"example.com/blog/internal/config"
	"example.com/blog/internal/entity"
	"example.com/blog/internal/modal"
	"example.com/blog/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	if err := service.RegisterUser(&user); err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	c.JSON(http.StatusOK, config.Success(nil))
}

func Login(c *gin.Context) {
	var user modal.UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	token, expAt, err := service.Login(&user)
	if err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	c.JSON(http.StatusOK, config.Success(gin.H{"token": token, "exp": expAt}))
}
