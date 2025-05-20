package main

import (
	"database/sql"
	"fmt"
)

// 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和
// transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
// 建表SQL脚本
// CREATE TABLE `accounts` (
// 	`id` bigint NOT NULL AUTO_INCREMENT,
// 	`balance` int DEFAULT 0,
// 	PRIMARY KEY (`id`)
//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
// CREATE TABLE `transactions` (
// 	`id` bigint NOT NULL AUTO_INCREMENT,
// 	`from_account_id` bigint,
// 	`to_account_id` bigint,
// 	`amount` int,
// 	PRIMARY KEY (`id`)
//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

// 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，
// 如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
func tranx(db *sql.DB) {
	accountA := 1
	accountB := 2

	// 开启事务
	tx, err0 := db.Begin()
	if err0 != nil {
		if tx != nil {
			tx.Rollback()
		}
		panic(err0)
	}
	// 查询账户A的余额
	var balanceA int
	err1 := tx.QueryRow("select balance from accounts where id = ?", accountA).Scan(&balanceA)
	if err1 != nil {
		tx.Rollback()
		panic(err1)
	}
	fmt.Printf("account A balance: %d\n", balanceA)
	// 校验余额
	if balanceA < 100 {
		tx.Rollback()
		panic("account A has insufficient balance")
	}

	//转账
	amount := 100
	_, err2 := tx.Exec("update accounts set balance = balance - ? where id = ?", amount, accountA)
	if err2 != nil {
		tx.Rollback()
		panic(err2)
	}
	_, err3 := tx.Exec("update accounts set balance = balance + ? where id = ?", amount, accountB)
	if err3 != nil {
		tx.Rollback()
		panic(err3)
	}
	_, err4 := tx.Exec("insert into transactions (from_account_id, to_account_id, amount) values (?, ?, ?)", accountA, accountB, amount)
	if err4 != nil {
		tx.Rollback()
		panic(err4)
	}
	tx.Commit()
}

// func main() {
// 	db := InitDB()
// 	tranx(db)
// }
