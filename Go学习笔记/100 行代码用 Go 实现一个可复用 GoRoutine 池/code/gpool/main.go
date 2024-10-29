package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	pool := NewPool(3)
	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		pool.SubmitTask(
			func() {
				defer wg.Done()
				time.Sleep(1 * time.Second)
				fmt.Println(i)
			})
	}

	wg.Wait()
	time.Sleep(3 * time.Second)
	pool.ReSize(1)

	for i := range 10 {
		wg.Add(1)
		pool.SubmitTask(
			func() {
				defer wg.Done()
				time.Sleep(1 * time.Second)
				fmt.Println(i)
			})
	}

	wg.Wait()

}
