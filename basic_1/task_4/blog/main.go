package main

import (
	"log/slog"

	"example.com/blog/internal/config"
	"example.com/blog/internal/entity"
	"example.com/blog/internal/router"
)

func main() {
	slog.Info("开始初始化数据库连接")
	db, err := config.InitMysql()
	if err != nil {
		slog.Error("初始化数据库连接失败", "error", err)
		return
	}
	slog.Info("初始化数据库连接成功")

	slog.Info("开始自动迁移模型完善数据库表")
	db.AutoMigrate(&entity.User{}, &entity.Post{}, &entity.Comment{})
	slog.Info("数据库表自动迁移完成")

	slog.Info("开始启动HTTP服务")
	router := router.SetupRouter()
	router.Run()
}
