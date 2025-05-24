package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局常量
const (
	USER_ID = "user_id"
)

// 全局DB变量
var GormDB *gorm.DB

// 配置文件结构体
type AppConfig struct {
	Database DatabaseConfig `yaml:"database"`
}

// 数据库连接信息结构体
type DatabaseConfig struct {
	DriverName string `yaml:"driver-name"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Name       string `yaml:"name"`
}

// 从配置文件读取数据库连接配置
func getDatabaseConfig() *DatabaseConfig {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	var config AppConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return &config.Database
}

// 初始化数据库连接
func InitMysql() (*gorm.DB, error) {
	// 获取数据库连接配置
	conf := getDatabaseConfig()
	// 初始化数据库连接
	dns := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, conf.Name)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// 获取底层DB连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 将数据库连接保存到全局变量中
	GormDB = db
	return db, sqlDB.Ping()
}
