package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
// 建表SQL脚本
// CREATE TABLE `students` (
//   `id` bigint NOT NULL AUTO_INCREMENT,
//   `name` varchar(20) DEFAULT '',
//   `age` int DEFAULT '0',
//   `grade` varchar(20) DEFAULT '',
//   PRIMARY KEY (`id`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

type Student struct {
	ID    uint
	Name  string
	Age   uint
	Grade string
}

// 连接数据库
func InitDB() *sql.DB {
	dsn := "admin:aa121212@tcp(localhost:3306)/gorm?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// 尝试与DB建立连接验证DSN是否正确
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

// 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"
func insert(db *sql.DB) {
	sqlInsert := "insert into students(name, age, grade) VALUES (?, ?, ?)"
	res, err0 := db.Exec(sqlInsert, "张三", 20, "三年级")
	if err0 != nil {
		panic(err0)
	}
	fmt.Printf("insert result: %+v\n", res)
}

// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息
func query(db *sql.DB) {
	sqlQuery := "select * from students where age > ?"
	rows, err1 := db.Query(sqlQuery, 18)
	if err1 != nil {
		panic(err1)
	}
	// 释放数据库连接
	defer rows.Close()
	// 遍历结果
	for rows.Next() {
		var student Student
		rows.Scan(&student.ID, &student.Name, &student.Age, &student.Grade)
		fmt.Printf("query result: %+v\n", student)
	}
}

// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"
func update(db *sql.DB) {
	sqlUpdate := "update students set grade = ? where name = ?"
	_, err := db.Exec(sqlUpdate, "四年级", "张三")
	if err != nil {
		panic(err)
	}
}

// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录
func del(db *sql.DB) {
	sqlDel := "delete from students where age < ?"
	_, err := db.Exec(sqlDel, 15)
	if err != nil {
		panic(err)
	}
}

// func main() {
// 	db := InitDB()
// 	// insert(db)
// 	// query(db)
// 	// update(db)
// 	del(db)
// }
