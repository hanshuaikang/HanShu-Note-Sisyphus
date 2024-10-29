package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type fn func()

type Pool struct {
	// 正在运行的 worker 数
	running int64
	// worker 列表
	workers []*Worker
	// worker 数量
	capacity int64
	// 锁
	lock sync.Mutex
}

func (p *Pool) incRunning() {
	atomic.AddInt64(&p.running, 1)
}

func (p *Pool) decRunning() {
	running := atomic.LoadInt64(&p.running)
	// 防止缩容的时候被打穿
	if running == 0 {
		return
	}
	atomic.AddInt64(&p.running, -1)
}

func (p *Pool) Running() int64 {
	return atomic.LoadInt64(&p.running)
}

func (p *Pool) Cap() int64 {
	return p.capacity
}

func NewPool(capacity int64) *Pool {
	return &Pool{
		capacity: capacity,
		workers:  make([]*Worker, 0),
		running:  0,
	}
}

func (p *Pool) PutWorker(worker *Worker) {
	p.lock.Lock()
	p.workers = append(p.workers, worker)
	p.decRunning()
	p.lock.Unlock()
}

func (p *Pool) SubmitTask(task fn) {
	worker := p.getWorker()
	p.incRunning()
	worker.task <- task
}

func (p *Pool) getWorker() *Worker {
	p.lock.Lock()
	defer p.lock.Unlock()

	if len(p.workers) > 0 {
		worker := p.workers[len(p.workers)-1]
		p.workers = p.workers[:len(p.workers)-1]
		return worker
	}

	if p.running < p.capacity {
		worker := &Worker{
			pool: p,
			task: make(chan fn),
		}
		worker.Run()
		return worker
	}

	for {
		p.lock.Unlock()
		time.Sleep(time.Millisecond)
		p.lock.Lock()

		if len(p.workers) > 0 {
			worker := p.workers[len(p.workers)-1]
			p.workers = p.workers[:len(p.workers)-1]
			return worker
		}
	}
}

// ReSize change the capacity of this pool
func (p *Pool) ReSize(size int64) {
	if size == p.Cap() {
		return
	}
	diff := p.capacity - size
	fmt.Println(fmt.Sprintf("diff: %d", len(p.workers)))
	atomic.StoreInt64(&p.capacity, size)
	if diff > 0 {
		for i := 0; i < int(diff); i++ {
			p.getWorker().task <- nil
		}
	}
}
