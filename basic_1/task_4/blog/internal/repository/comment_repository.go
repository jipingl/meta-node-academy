package repository

import (
	"errors"

	"example.com/blog/internal/config"
	"example.com/blog/internal/entity"
)

func CreateComment(comment *entity.Comment) (uint, error) {
	res := config.GormDB.Create(comment)
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, errors.New("create comment failed")

	}
	return comment.ID, nil
}

func GetCommentList(postID uint) ([]entity.Comment, error) {
	var comments []entity.Comment
	res := config.GormDB.Debug().Where("post_id = ?", postID).Preload("User").Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("comment not found")
	}
	return comments, nil
}
