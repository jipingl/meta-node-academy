package router

import (
	"net/http"

	"example.com/blog/internal/controller"
	"example.com/blog/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 测试地址
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 注册用户
	router.POST("/register", controller.RegisterUser)
	// 登录
	router.POST("/login", controller.Login)

	//需要认证的api
	apiRouter := router.Group("/api")
	apiRouter.Use(middleware.AuthMiddleware())
	{
		// 创建文章
		apiRouter.POST("/post", controller.CreatePost)
		// 获取文章列表
		apiRouter.GET("/post", controller.GetPostList)
		// 根据文章ID获取文章详情
		apiRouter.GET("/post/:id", controller.GetPostDetail)
		// 更新文章
		apiRouter.PUT("/post/:id", controller.UpdatePost)
		// 删除文章
		apiRouter.DELETE("/post/:id", controller.DeletePost)
		// 发表评论
		apiRouter.POST("/comment", controller.CreateComment)
		// 获取文章评论列表
		apiRouter.GET("/comment", controller.GetCommentList)
	}

	return router
}
