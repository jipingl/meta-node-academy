package controller

import (
	"errors"
	"net/http"
	"strconv"

	"example.com/blog/internal/config"
	"example.com/blog/internal/entity"
	"example.com/blog/internal/modal"
	"example.com/blog/internal/service"
	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	var commentCreate modal.CommentCreate
	if err := c.ShouldBindJSON(&commentCreate); err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	var comment entity.Comment
	comment.Content = commentCreate.Content
	comment.PostID = commentCreate.PostID
	comment.UserID = c.MustGet(config.USER_ID).(uint)
	commentID, err := service.CreateComment(&comment)
	if err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	c.JSON(http.StatusOK, config.Success(commentID))
}

func GetCommentList(c *gin.Context) {
	// 获取文章ID
	postIDStr := c.Query("post_id")
	// 验证文章ID
	if postIDStr == "" {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, errors.New("post_id is required")))
		return
	}
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	comments, err := service.GetCommentList(uint(postID))
	if err != nil {
		c.JSON(http.StatusOK, config.Fail(http.StatusInternalServerError, err))
		return
	}
	// 将实体转换为响应结构体
	var commentList []modal.CommentView
	for _, comment := range comments {
		commentList = append(commentList, modal.CommentView{
			ID:        comment.ID,
			Content:   comment.Content,
			PostID:    comment.PostID,
			UserID:    comment.UserID,
			Username:  comment.User.Username,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, config.Success(commentList))
}
