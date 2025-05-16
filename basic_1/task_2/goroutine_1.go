package main

import (
	"fmt"
	"time"
)

// 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
// 考察点 ： go 关键字的使用、协程的并发执行。
func topic3() {
	// 打印1到10的奇数
	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 1 {
				fmt.Printf("奇数协程：%d\n", i)
				time.Sleep(time.Millisecond)
			}
		}
	}()
	// 打印2到10的偶数
	go func() {
		for i := 2; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Printf("偶数协程：%d\n", i)
				time.Sleep(time.Millisecond)
			}
		}
	}()
	time.Sleep(time.Second)
}

// func main() {
// 	topic3()
// }
