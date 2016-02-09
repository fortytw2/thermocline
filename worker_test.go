package thermocline_test

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bradfitz/iter"
	"github.com/fortytw2/thermocline"
	"github.com/fortytw2/thermocline/brokers/mem"
)

func TestWorker(t *testing.T) {
	t.Parallel()

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

	tn := rand.Intn(256)
	//TODO(fortytw2): don't use iter for tests, even though it's nice
	for i := range iter.N(tn) {
		task, err := thermocline.NewTask(fmt.Sprintf("test %d", i))
		if err != nil {
			t.Error("could not create test task", err)
		}

		writer <- task
	}

	stopper := make(chan struct{})
	var worked int64
	wg := &sync.WaitGroup{}
	for _ = range iter.N(rand.Intn(rand.Intn(256))) {
		wg.Add(1)
		go thermocline.NewWorker(reader, writer, func(task *thermocline.Task) ([]*thermocline.Task, error) {
			atomic.AddInt64(&worked, 1)
			return nil, nil
		}, stopper).Work(wg)
	}

	time.Sleep(500 * time.Millisecond)
	close(stopper)

	wg.Wait()

	if atomic.LoadInt64(&worked) != int64(tn) {
		t.Errorf("%d tasks not worked in basic test after 500ms, instead %d", tn, atomic.LoadInt64(&worked))
	}
}

func TestWorkerRetries(t *testing.T) {
	t.Parallel()

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

	tn := rand.Intn(256)
	for i := range iter.N(tn) {
		task, err := thermocline.NewTask(fmt.Sprintf("test %d", i))
		if err != nil {
			t.Error("could not create test task", err)
		}

		writer <- task
	}

	stopper := make(chan struct{})
	var worked int64
	wg := &sync.WaitGroup{}
	for _ = range iter.N(rand.Intn(256)) {
		wg.Add(1)
		go thermocline.NewWorker(reader, writer, func(task *thermocline.Task) ([]*thermocline.Task, error) {
			atomic.AddInt64(&worked, 1)
			return nil, errors.New("cannot process task, herp derup")
		}, stopper).Work(wg)
	}

	time.Sleep(500 * time.Millisecond)
	close(stopper)

	wg.Wait()

	if atomic.LoadInt64(&worked) != int64(tn*3) {
		t.Errorf("%d tasks not worked in retry test after 500ms, actually %d", tn*3, atomic.LoadInt64(&worked))
	}
}

func TestWorkerStop(t *testing.T) {

}
