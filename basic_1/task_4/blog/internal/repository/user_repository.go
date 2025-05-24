package repository

import (
	"errors"

	"example.com/blog/internal/config"
	"example.com/blog/internal/entity"
)

func CreateUser(user *entity.User) error {
	res := config.GormDB.Create(user)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("failed to create user")
	}
	return nil
}

func GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	res := config.GormDB.Where("username = ?", username).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("username or password error")
	}
	return &user, nil
}
