package mem

import (
	"fmt"
	"sync"

	"github.com/fortytw2/thermocline"
)

// Broker is a pure in memory message broker
type Broker struct {
	queues   map[string]chan *thermocline.Task
	failures map[string]int
	*sync.Mutex
}

// NewBroker returns the allocated message broker
func NewBroker() *Broker {
	return &Broker{
		queues: make(map[string]chan *thermocline.Task),
		Mutex:  &sync.Mutex{},
	}
}

// Open returns the broker
func (b *Broker) Open(queue string, version string) (chan *thermocline.Task, error) {
	key := fmt.Sprintf("%s:%s", queue, version)

	c, ok := b.queues[key]
	if !ok {
		b.queues[key] = make(chan *thermocline.Task, 1024)
		return b.queues[key], nil
	}

	return c, nil
}

func (b *Broker) Failures(queue, version string) (int, error) {
	key := fmt.Sprintf("%s:%s", queue, version)
	b.Lock()
	defer b.Unlock()
	f, ok := b.failures[key]
	if !ok {
		b.failures[key] = 0
		return 0, nil
	}
	return f, nil
}
