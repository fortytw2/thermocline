package thermocline

import "sync"

// A Worker is the basic unit of task execution
type Worker struct {
	reader    <-chan *Task
	writer    chan<- *Task
	processFn Processor
	stop      chan struct{}
}

// NewWorker creates a new worker
func NewWorker(reader <-chan *Task, writer chan<- *Task, fn Processor, stopper chan struct{}) *Worker {
	return &Worker{
		reader:    reader,
		writer:    writer,
		processFn: fn,
		stop:      stopper,
	}
}

// Work turns on the worker
func (w *Worker) Work(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		// safely stop the worker
		case <-w.stop:
			return
		case task := <-w.reader:
			tasks, err := w.processFn(task)
			if err != nil {
				if task.Retries < MaxRetries-1 {
					task.Retries++
					w.writer <- task
					continue
				}
			}

			// submit any new tasks returned by the old one
			if tasks != nil {
				for _, t := range tasks {
					w.writer <- t
				}
			}
		}
	}
}
