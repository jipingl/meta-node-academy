package repository

import (
	"errors"

	"example.com/blog/internal/config"
	"example.com/blog/internal/entity"
)

func CreatePost(post *entity.Post) (uint, error) {
	res := config.GormDB.Create(post)
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, errors.New("failed to create post")
	}
	return post.ID, nil
}

func GetPostList() ([]entity.Post, error) {
	var posts []entity.Post
	res := config.GormDB.Order("id desc").Find(&posts)
	if res.Error != nil {
		return nil, res.Error
	}
	return posts, nil
}

func GetPostDetail(id uint) (entity.Post, error) {
	var post entity.Post
	res := config.GormDB.First(&post, id)
	if res.Error != nil {
		return entity.Post{}, res.Error
	}
	if res.RowsAffected == 0 {
		return entity.Post{}, errors.New("post not found")
	}
	return post, nil
}

func UpdatePost(post *entity.Post) error {
	res := config.GormDB.Save(post)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("post not found")
	}
	return nil
}

func DeletePost(id uint) error {
	res := config.GormDB.Delete(&entity.Post{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("post not found")
	}
	return nil
}
