package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary
// CREATE TABLE `employees` (
//   `id` bigint NOT NULL AUTO_INCREMENT,
//   `name` varchar(20) DEFAULT NULL,
//   `department` varchar(20) DEFAULT NULL,
//   `salary` int DEFAULT NULL,
//   PRIMARY KEY (`id`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

type Employee struct {
	ID         uint
	Name       string
	Department string
	Salary     uint
}

func SqlxInitDB() *sqlx.DB {
	dsn := "admin:aa121212@tcp(localhost:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// 测试连接
	err1 := db.Ping()
	if err1 != nil {
		panic(err1)
	}
	return db
}

// 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。,
func query_1(db *sqlx.DB) {
	var employees []Employee
	err := db.Select(&employees, "select * from employees where department = ?", "技术部")
	if err != nil {
		panic(err)
	}
	fmt.Printf("all employees: %+v\n", employees)
}

// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中
func query_2(db *sqlx.DB) {
	var employee Employee
	err := db.Get(&employee, "select * from employees order by salary desc limit 1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("the highest salary employee: %+v\n", employee)
}

// func main() {
// 	db := SqlxInitDB()
// 	query_2(db)
// }
