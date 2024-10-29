package main

type Worker struct {
	task chan fn
	pool *Pool
}

func (w *Worker) Run() {
	go func() {
		for t := range w.task {
			// 退出机制，本文不会涉及
			if t == nil {
				// 正在执行的数量 -1
				w.pool.decRunning()
				return
			}
			t()
			// 执行完了将该 worker 重新放回池子里
			w.pool.PutWorker(w)
		}
	}()
}
