package controller

import (
	"net/http"
	"strconv"

	"example.com/blog/internal/config"
	"example.com/blog/internal/entity"
	"example.com/blog/internal/modal"
	"example.com/blog/internal/service"
	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	var postCreate modal.PostCreate
	if err := c.ShouldBindJSON(&postCreate); err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	var post entity.Post
	post.Title = postCreate.Title
	post.Content = postCreate.Content
	post.UserID = c.MustGet(config.USER_ID).(uint)
	postID, err := service.CreatePost(&post)
	if err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	c.JSON(http.StatusOK, config.Success(postID))
}

func GetPostList(c *gin.Context) {
	posts, err := service.GetPostList()
	if err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	// 将实体转换为响应结构体
	var postList []modal.PostView
	for _, post := range posts {
		postList = append(postList, modal.PostView{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			UserID:    post.UserID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, config.Success(postList))
}

func GetPostDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, e := strconv.ParseUint(idStr, 10, 64)
	if e != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, e))
		return
	}
	post, err := service.GetPostDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	// 将实体转换为响应结构体
	postView := modal.PostView{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
	c.JSON(http.StatusOK, config.Success(postView))
}

func UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, e := strconv.ParseUint(idStr, 10, 64)
	if e != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, e))
		return
	}
	var postUpdate modal.PostUpdate
	if err := c.ShouldBindJSON(&postUpdate); err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	// 当前登录的用户ID
	userID := c.MustGet(config.USER_ID).(uint)
	err := service.UpdatePost(uint(id), &postUpdate, userID)
	if err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	c.JSON(http.StatusOK, config.Success(nil))
}

func DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, e := strconv.ParseUint(idStr, 10, 64)
	if e != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, e))
		return
	}
	// 当前登录的用户ID
	userID := c.MustGet(config.USER_ID).(uint)
	err := service.DeletePost(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	c.JSON(http.StatusOK, config.Success(nil))
}
