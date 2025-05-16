package main

import (
	"fmt"
	"sync"
	"time"
)

// 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// 考察点 ：协程原理、并发任务调度。

// 定义任务类型
type Task struct {
	ID   int           // 任务ID
	Fn   func()        // 要执行的具体任务
	Res  interface{}   // 执行结果
	Err  error         // 错误信息
	Cost time.Duration // 执行耗时
}

// 定义任务调度器类型
type Scheduler struct {
	tasks       []*Task        // 任务列表
	concurrency int            // 并发度
	wg          sync.WaitGroup // 等待组
	startTime   time.Time      // 开始时间
}

// 创建调度器的函数
func NewScheduler(concurrency int) *Scheduler {
	return &Scheduler{concurrency: concurrency}
}

// 向调度器中添加任务
func (s *Scheduler) AddTask(fn func()) {
	task := Task{
		ID: len(s.tasks) + 1,
		Fn: fn,
	}
	s.tasks = append(s.tasks, &task)
}

// 工做协程
func (s *Scheduler) worker(taskChan chan *Task) {
	// 获取任务并执行
	for taskPtr := range taskChan {
		start := time.Now()

		// 执行任务
		func() {
			defer func() {
				// 捕获执行过程中出现的异常
				if r := recover(); r != nil {
					taskPtr.Err = fmt.Errorf("panic: %v", r)
				}
			}()
			// 执行具体的任务
			taskPtr.Fn()
		}()

		taskPtr.Cost = time.Since(start)
		s.wg.Done()
	}
}

// 执行所有任务
func (s *Scheduler) Run() {
	s.startTime = time.Now()
	taskChan := make(chan *Task, len(s.tasks))

	// 启动多个 worker
	for i := 0; i < s.concurrency; i++ {
		go s.worker(taskChan)
	}

	// 开始分发任务（将任务写入chan）
	s.wg.Add(len(s.tasks))
	for _, task := range s.tasks {
		taskChan <- task
	}
	close(taskChan)
	// 阻塞等待任务执行完毕
	s.wg.Wait()
}

// 打印任务统计信息
func (s *Scheduler) PrintStats() {
	fmt.Printf("\n=== 任务统计 ===\n")
	fmt.Printf("总任务数：%d\n", len(s.tasks))
	fmt.Printf("并发数：%d\n", s.concurrency)
	fmt.Printf("总耗时：%v\n", time.Since(s.startTime))

	fmt.Printf("\n=== 任务详情 ===\n")
	for _, task := range s.tasks {
		status := "成功"
		if task.Err != nil {
			status = fmt.Sprintf("任务执行失败：%v", task.Err)
		}
		fmt.Printf("任务 %d - 耗时：%v - 状态：%s\n", task.ID, task.Cost, status)
	}
}

// func main() {
// 	// 新建任务调度器
// 	scheduler := NewScheduler(3)

// 	// 添加示例任务
// 	scheduler.AddTask(func() {
// 		time.Sleep(1 * time.Second)
// 		fmt.Println("任务1执行完成")
// 	})

// 	scheduler.AddTask(func() {
// 		time.Sleep(2 * time.Second)
// 		fmt.Println("任务2执行完成")
// 	})

// 	scheduler.AddTask(func() {
// 		time.Sleep(500 * time.Millisecond)
// 		fmt.Println("任务3执行完成")
// 	})

// 	scheduler.AddTask(func() {
// 		time.Sleep(1 * time.Second)
// 		panic("任务4出错了")
// 	})

// 	// 执行任务
// 	scheduler.Run()

// 	// 打印统计信息
// 	scheduler.PrintStats()
// }
