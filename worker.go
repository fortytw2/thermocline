package thermocline

import "sync"

// A Worker is the basic unit of task execution
type Worker struct {
	ingress   <-chan *Task
	egress    chan<- *Task
	processFn Processor
	stop      chan struct{}
}

// NewWorker creates a new worker
func NewWorker(ingress <-chan *Task, egress chan<- *Task, fn Processor, stopper chan struct{}) *Worker {
	return &Worker{
		ingress:   ingress,
		egress:    egress,
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
		case task := <-w.ingress:
			tasks, err := w.processFn(task)
			if err != nil {
				if task.Retries < MaxRetries {
					task.Retries++
					w.egress <- task
					continue
				}
			}

			// submit any new tasks returned by the old one
			if tasks != nil {
				for _, t := range tasks {
					w.egress <- t
				}
			}
		}
	}
}
