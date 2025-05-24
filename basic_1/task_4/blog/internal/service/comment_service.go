package service

import (
	"example.com/blog/internal/entity"
	"example.com/blog/internal/repository"
)

func CreateComment(comment *entity.Comment) (uint, error) {
	id, err := repository.CreateComment(comment)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetCommentList(postID uint) ([]entity.Comment, error) {
	comments, err := repository.GetCommentList(postID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
