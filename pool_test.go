package thermocline_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/fortytw2/thermocline"
	"github.com/fortytw2/thermocline/brokers/mem"
)

func TestPool(t *testing.T) {
	t.Parallel()

	var b thermocline.Broker
	b = mem.NewBroker()

	ticker := time.NewTicker(time.Millisecond * 10)
	go func() {
		w, err := b.Write("test", thermocline.NoVersion)
		if err != nil {
			t.Fatalf("cannot get write chan %s", err)
		}
		for {
			select {
			case t := <-ticker.C:
				tk, _ := thermocline.NewTask(t)
				w <- tk
			}
		}
	}()

	var worked int64
	p, err := thermocline.NewPool("test", thermocline.NoVersion, b, func(task *thermocline.Task) ([]*thermocline.Task, error) {
		atomic.AddInt64(&worked, 1)
		return nil, nil
	}, 5)
	if err != nil {
		t.Errorf("cannot create pool %s", err)
	}

	time.Sleep(1 * time.Second)
	err = p.Add(120)
	if err != nil {
		t.Errorf("cannot add workers, %s", err)
	}

	time.Sleep(1 * time.Second)
	err = p.Add(-37)
	if err != nil {
		t.Errorf("cannot remove workers, %s", err)
	}

	time.Sleep(1 * time.Second)
	ticker.Stop()

	err = p.Stop()
	if err != nil {
		t.Errorf("cannot stop pool, %s", err)
	}

	if p.Len() != 0 {
		t.Errorf("pool length is not 0! - %d", p.Len())
	}

	if a := atomic.LoadInt64(&worked); a <= 300 {
		t.Errorf("more than 300 tasks should be worked, %d", a)
	}
}
