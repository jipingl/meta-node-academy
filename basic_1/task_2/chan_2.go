package main

import (
	"fmt"
	"time"
)

// 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
// 考察点 ：通道的缓冲机制。
func topic5() {
	// 创建一个带有缓冲的通道
	ch := make(chan int, 10)

	// 向通道发送100个整数
	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i
		}
		close(ch)
		fmt.Printf("Send Finished\n")
	}()

	// 接收协程
	go func() {
		for num := range ch {
			fmt.Printf("Recieved: %d\n", num)
		}
	}()

	// 主程等待
	time.Sleep(time.Second * 5)
}

// func main() {
// 	topic5()
// }
