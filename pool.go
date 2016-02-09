package thermocline

import (
	"fmt"
	"sync"

	"github.com/bradfitz/iter"
)

// A Pool is a set of workers that all function on the same queue
type Pool struct {
	Queue   string
	Version string

	read  <-chan *Task
	write chan<- *Task
	fn    Processor
	// workers are only kept track of by their stop channel
	workers []chan struct{}
	wg      *sync.WaitGroup
	*sync.RWMutex
}

// NewPool returns a running worker pool on the given queue/version
func NewPool(queue, version string, b Broker, fn Processor, workers int) (*Pool, error) {
	r, err := b.Read(queue, version)
	if err != nil {
		return nil, err
	}

	w, err := b.Write(queue, version)
	if err != nil {
		return nil, err
	}

	p := &Pool{
		Queue:   queue,
		Version: version,

		read:    r,
		write:   w,
		fn:      fn,
		workers: make([]chan struct{}, 0),

		wg:      &sync.WaitGroup{},
		RWMutex: &sync.RWMutex{},
	}

	err = p.Add(workers)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Add changes the number of workers
func (p *Pool) Add(n int) error {
	if n >= 1 {
		for range iter.N(n) {
			stopC := make(chan struct{})
			p.wg.Add(1)
			go NewWorker(p.read, p.write, p.fn, stopC).Work(p.wg)
			p.Lock()
			p.workers = append(p.workers, stopC)
			p.Unlock()
		}
	} else if n <= -1 {
		for range iter.N(-1 * n) {
			p.Lock()
			// close channel to stop worker
			close(p.workers[len(p.workers)-1])

			// from github.com/golang/go/wiki/SliceTricks
			// Delete without preserving order
			p.workers[len(p.workers)-1] = nil
			p.workers = p.workers[:len(p.workers)-1]
			// unlock
			p.Unlock()
		}
	} else {
		return fmt.Errorf("%d is not a valid number of workers to add", n)
	}
	return nil
}

// Len returns the total number of workers in this group
func (p *Pool) Len() int {
	p.RLock()
	defer p.RUnlock()
	return len(p.workers)
}

// Stop turns off all workers in the pool
func (p *Pool) Stop() error {
	err := p.Add(-1 * p.Len())
	if err != nil {
		return err
	}
	p.wg.Wait()
	return nil
}
