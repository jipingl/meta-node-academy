package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ：原子操作、并发数据安全。
func topic6() {
	// 声明一个原子计数器
	var counter atomic.Uint32
	// 创建协程
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				counter.Add(1)
			}
		}()
	}
	// 主程等待
	time.Sleep(time.Second)
	// 打印结果
	fmt.Printf("Result: %d\n", counter.Load())
}

// func main() {
// 	topic6()
// }
