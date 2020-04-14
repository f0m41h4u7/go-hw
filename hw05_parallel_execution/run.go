package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks
func Run(tasks []Task, n int, m int) error {
	var (
		wg        sync.WaitGroup
		errNumber int
		mt        sync.Mutex
	)
	wg.Add(n)

	ch := make(chan Task, len(tasks))
	for _, t := range tasks {
		ch <- t
	}
	close(ch)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range ch {
				err := task()
				mt.Lock()
				if err != nil {
					errNumber++
				}
				if errNumber == m {
					mt.Unlock()
					return
				}
				mt.Unlock()
			}
		}()
	}
	wg.Wait()
	if errNumber >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
