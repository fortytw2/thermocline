package thermocline_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bradfitz/iter"
	"github.com/fortytw2/thermocline"
	"github.com/fortytw2/thermocline/brokers/mem"
)

func TestWorker(t *testing.T) {
	var broker thermocline.Broker
	broker = mem.NewBroker()

	reader, err := broker.Read("test", thermocline.NoVersion)
	if err != nil {
		t.Errorf("could not open queue '%s'", err)
	}

	writer, err := broker.Write("test", thermocline.NoVersion)
	if err != nil {
		t.Errorf("could not open queue '%s'", err)
	}

	for i := range iter.N(128) {
		task, err := thermocline.NewTask(fmt.Sprintf("test %d", i))
		if err != nil {
			t.Error("could not create test task", err)
		}

		writer <- task
	}

	stopper := make(chan struct{})
	var worked int64
	wg := &sync.WaitGroup{}
	for _ = range iter.N(100) {
		wg.Add(1)
		go thermocline.NewWorker(reader, writer, func(task *thermocline.Task) ([]*thermocline.Task, error) {
			atomic.AddInt64(&worked, 1)
			return nil, nil
		}, stopper).Work(wg)
	}

	time.Sleep(500 * time.Millisecond)
	close(stopper)

	wg.Wait()

	if atomic.LoadInt64(&worked) != 128 {
		t.Error("128 tasks not worked in basic test after 500ms")
	}
}

func TestWorkerRetries(t *testing.T) {
	// 	broker := mem.NewBroker()
	//
	// 	in, err := broker.Read("test", thermocline.NoVersion)
	// 	if err != nil {
	// 		t.Errorf("could not open queue '%s'", err)
	// 	}
	//
	// 	task, err := thermocline.NewTask("test")
	// 	if err != nil {
	// 		t.Error("could not create test task", err)
	// 	}
	//
	// 	stopper := make(chan struct{})
	// 	w := thermocline.NewWorker(in, func(task *thermocline.Task) ([]*thermocline.Task, error) {
	// 		return nil, errors.New("cannot proccess task")
	// 	}, stopper)
	//
	// 	go w.Work()
	// 	in <- task

}

func TestWorkerStop(t *testing.T) {

}
