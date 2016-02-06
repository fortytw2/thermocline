package thermocline

// A Worker is the basic unit of task execution
type Worker struct {
	queue     chan *Task
	processFn Processor
	stop      chan struct{}
}

func NewWorker(queue chan *Task, fn Processor, stopper chan struct{}) *Worker {
	return &Worker{
		queue:     queue,
		processFn: fn,
		stop:      stopper,
	}
}

// Work turns on the worker
func (w *Worker) Work() {
	for {
		select {
		// safely stop the worker
		case <-w.stop:
			return
		case task := <-w.queue:
			tasks, err := w.processFn(task)
			if err != nil {
				if task.Retries < MaxRetries {
					task.Retries++
					w.queue <- task
					continue
				}
			}

			// submit any new tasks returned by the old one
			if tasks != nil {
				for _, t := range tasks {
					w.queue <- t
				}
			}
		}
	}
}
