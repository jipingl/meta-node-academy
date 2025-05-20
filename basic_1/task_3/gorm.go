package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 题目1：模型定义,
// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
// 要求 ：
// 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。,
type User struct {
	gorm.Model
	Name      string
	PostCount uint `gorm:"default:0"`
	Posts     []Post
}

type Post struct {
	gorm.Model
	Title         string
	Content       string
	UserID        uint
	CommentStatus string `gorm:"default:no_comments"`
	Comments      []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
}

func GormInitDB() *gorm.DB {
	dsn := "admin:aa121212@tcp(localhost:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 设置连接池参数
	sqlDB, err1 := db.DB()
	if err1 != nil {
		panic(err1)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}

// 编写Go代码，使用Gorm创建这些模型对应的数据库表。
func createTables(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic(err)
	}
}

func dataReady(db *gorm.DB) {
	users := []User{
		User{
			Name: "AA",
			Posts: []Post{
				Post{
					Title:   "AA-P0",
					Content: "AA-P0",
					Comments: []Comment{
						Comment{
							Content: "AA-P0-C0",
						},
						Comment{
							Content: "AA-P0-C1",
						},
					},
				},
				Post{
					Title:   "AA-P1",
					Content: "AA-P1",
					Comments: []Comment{
						Comment{
							Content: "AA-P1-C0",
						},
						Comment{
							Content: "AA-P1-C1",
						},
						Comment{
							Content: "AA-P1-C2",
						},
					},
				},
			},
		},
		User{
			Name: "BB",
			Posts: []Post{
				Post{
					Title:   "BB-P0",
					Content: "BB-P0",
					Comments: []Comment{
						Comment{
							Content: "BB-P0-C0",
						},
					},
				},
			},
		},
	}
	db.Create(&users)
}

// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
func queryUserPostComment(db *gorm.DB) {
	var users []User
	db.Debug().Preload("Posts.Comments").Find(&users, 7)
	fmt.Printf("%v\n", users)
}

// 编写Go代码，使用Gorm查询评论数量最多的文章信息。
func postWithMostComment(db *gorm.DB) {
	var post Post
	subQuery := db.Select("post_id").Group("post_id").Order("count(id) desc").Limit(1).Table("comments")
	db.Where("id=(?)", subQuery).First(&post)
	fmt.Printf("%v\n", post)
}

// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	rtx := tx.Model(&User{}).Where("id=?", p.UserID).UpdateColumn("post_count", gorm.Expr("post_count+?", 1))
	if rtx.Error != nil {
		return rtx.Error
	}
	return nil
}

// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	tx.Unscoped().Where("id=?", c.ID).First(&c)
	// 查询文章的有效评论数量
	var count int64
	tx.Debug().Model(&Comment{}).Where("post_id=?", c.PostID).Count(&count)
	if count == 0 {
		tx.Debug().Model(&Post{}).Where("id=?", c.PostID).Updates(Post{CommentStatus: "no_comments"})
	}
	return nil
}

// func main() {
// 	db := GormInitDB()
// 	// createTables(db)
// 	// dataReady(db)
// 	// queryUserPostComment(db)
// 	// postWithMostComment(db)

// 	// 创建新的文章
// 	// db.Create(&Post{Title: "BB-P2", Content: "BB-P2", UserID: 7})

// 	// 删除文章的唯一一个评论
// 	db.Delete(&Comment{Model: gorm.Model{ID: 6}})
// }
