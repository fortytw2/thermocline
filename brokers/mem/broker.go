package mem

import (
	"fmt"
	"sync"

	"github.com/fortytw2/thermocline"
)

// Broker is a pure in memory message broker
type Broker struct {
	egress  map[string]chan *thermocline.Task
	ingress map[string]chan *thermocline.Task

	stats *thermocline.Stats
	*sync.Mutex
}

// NewBroker returns the allocated message broker
func NewBroker() *Broker {
	b := &Broker{
		egress:  make(map[string]chan *thermocline.Task),
		ingress: make(map[string]chan *thermocline.Task),
		stats: &thermocline.Stats{
			Total: make(map[string]int),
		},
		Mutex: &sync.Mutex{},
	}
	return b
}

func (b *Broker) monitor(queue, version string) {
	key := fmt.Sprintf("%s:%s", queue, version)
	for {
		select {
		case t := <-b.ingress[key]:
			if t.Retries == 0 {
				b.Lock()
				b.stats.Total[key]++
				b.Unlock()
				b.egress[key] <- t
			} else {
				b.egress[key] <- t
			}
		}
	}
}

// Open returns the broker
func (b *Broker) Read(queue string, version string) (<-chan *thermocline.Task, error) {
	key := fmt.Sprintf("%s:%s", queue, version)

	b.Lock()
	defer b.Unlock()
	c, ok := b.egress[key]
	if !ok {
		b.egress[key] = make(chan *thermocline.Task, 1024)
		return b.egress[key], nil
	}

	return c, nil
}

// Open returns the broker
func (b *Broker) Write(queue string, version string) (chan<- *thermocline.Task, error) {
	key := fmt.Sprintf("%s:%s", queue, version)

	b.Lock()
	defer b.Unlock()
	c, ok := b.ingress[key]
	if !ok {
		b.ingress[key] = make(chan *thermocline.Task, 1024)
		go b.monitor(queue, version)
		return b.ingress[key], nil
	}

	return c, nil
}

// Stats returns all queue stats
func (b *Broker) Stats() *thermocline.Stats {
	b.Lock()
	defer b.Unlock()

	out := *b.stats
	return &out
}
