package service

import (
	"errors"

	"example.com/blog/internal/entity"
	"example.com/blog/internal/modal"
	"example.com/blog/internal/repository"
)

func CreatePost(post *entity.Post) (uint, error) {
	return repository.CreatePost(post)
}

func GetPostList() ([]entity.Post, error) {
	return repository.GetPostList()
}

func GetPostDetail(id uint) (entity.Post, error) {
	return repository.GetPostDetail(id)
}

func UpdatePost(id uint, postUpdate *modal.PostUpdate, userID uint) error {
	post, err := repository.GetPostDetail(id)
	if err != nil {
		return err
	}
	// 校验是否是文章作者
	if post.UserID != userID {
		return errors.New("not allowed")
	}
	post.Title = postUpdate.Title
	post.Content = postUpdate.Content
	return repository.UpdatePost(&post)
}

func DeletePost(id uint, userID uint) error {
	post, err := repository.GetPostDetail(id)
	if err != nil {
		return err
	}
	// 校验是否是文章作者
	if post.UserID != userID {
		return errors.New("not allowed")
	}
	return repository.DeletePost(id)
}
