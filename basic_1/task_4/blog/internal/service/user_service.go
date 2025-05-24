package service

import (
	"errors"

	"example.com/blog/internal/entity"
	"example.com/blog/internal/modal"
	"example.com/blog/internal/repository"
	"example.com/blog/internal/util"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(user *entity.User) error {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)
	return repository.CreateUser(user)
}

func Login(user *modal.UserLogin) (token string, expAt int64, err error) {
	// 查询用户
	storeUser, err := repository.GetUserByUsername(user.Username)
	if err != nil {
		return "", 0, errors.New("username or password error")
	}
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(storeUser.Password), []byte(user.Password))
	if err != nil {
		return "", 0, errors.New("username or password error")
	}
	// 生成jwt
	return util.GenerateToken(storeUser)
}
