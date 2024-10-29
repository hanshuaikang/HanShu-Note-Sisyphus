package main

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	AntsSize = 1000
	n        = 1000000
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
)

var curMem uint64

func demoFunc() {
	time.Sleep(time.Duration(10) * time.Millisecond)
}

func TestACurMem(t *testing.T) {
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	cur := mem.TotalAlloc / KiB
	t.Logf("start memory usage:%d KiB", cur)

}

// 基准测试函数
func TestPoolWaitToGetWorker(t *testing.T) {
	var wg sync.WaitGroup
	p := NewPool(AntsSize)

	for i := 0; i < n; i++ {
		wg.Add(1)
		p.SubmitTask(func() {
			demoFunc()
			wg.Done()
		})
	}
	wg.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("poll: memory usage:%d MB", curMem)
}

func TestNoPool(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			demoFunc()
			wg.Done()
		}()
	}

	wg.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("no poll: memory usage:%d MB", curMem)
}
