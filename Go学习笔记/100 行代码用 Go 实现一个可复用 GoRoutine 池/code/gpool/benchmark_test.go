package main

import (
	"sync"
	"testing"
)

// 基准测试函数
func BenchmarkWorkerPool(b *testing.B) {
	pool := NewPool(1000) // 设置池的容量
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		pool.SubmitTask(func() {
			// 模拟任务执行
			defer wg.Done()
			demoFunc()
		})
	}
	wg.Wait()
}

// 原生 goroutine 的基准测试
func BenchmarkNativeGoroutine(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 模拟任务执行
			demoFunc()
		}()
	}
	wg.Wait()
}
