package worker

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Worker interface {
	Start(ctx context.Context)
	Stop()
}

type settings struct {
	interval             time.Duration
	concurrentTasksCount int
}

func NewWorker(
	ctx context.Context,
	task func(ctx context.Context),
	interval time.Duration,
	concurrentTasksCount int,
) Worker {
	return &worker{
		ctx:  ctx,
		task: task,
		settings: &settings{
			interval:             interval,
			concurrentTasksCount: concurrentTasksCount,
		},
	}
}

type worker struct {
	settings *settings

	ctx          context.Context
	task         func(ctx context.Context)
	interval     time.Duration
	concurrentCh chan struct{}
	stopCh       chan struct{}
	stopChMutex  sync.Mutex
	taskWg       sync.WaitGroup
}

func (w *worker) Start(context.Context) {
	if w.stopCh != nil {
		w.Stop()
	}

	w.concurrentCh = make(chan struct{}, w.settings.concurrentTasksCount)
	w.stopCh = make(chan struct{})

	go func() {
		for {
			select {
			case <-w.ctx.Done():
				return
			case <-w.stopCh:
				return
			case <-time.After(w.settings.interval):
				if len(w.concurrentCh) < cap(w.concurrentCh) {
					w.concurrentCh <- struct{}{}
					w.taskWg.Add(1)
					go w.runTask()
				}
			}
		}
	}()
}

func (w *worker) runTask() {
	defer func() {
		<-w.concurrentCh
		w.taskWg.Done()
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered", r)
		}
	}()

	w.task(w.ctx)
}

func (w *worker) Stop() {
	w.stopChMutex.Lock()
	defer w.stopChMutex.Unlock()

	if w.stopCh == nil {
		return // Работник уже остановлен.
	}

	w.stopCh <- struct{}{}
	close(w.stopCh)
	w.taskWg.Wait()

	w.stopCh = nil
}
