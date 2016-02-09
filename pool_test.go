package thermocline_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/fortytw2/thermocline"
)

func TestPool(t *testing.T) {
	t.Parallel()

	testPool(t, 0)
}

func TestPoolSlowWorkers(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping in short mode.")
	}
	testPool(t, 250*time.Millisecond)
}

func testPool(t *testing.T, duration time.Duration) {
	b, _, writer := SetupBroker(t)
	ticker := time.NewTicker(time.Millisecond * 10)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				tk, _ := thermocline.NewTask(t)
				writer <- tk
			}
		}
	}()

	var worked int64
	p, err := thermocline.NewPool("test", thermocline.NoVersion, b, func(task *thermocline.Task) ([]*thermocline.Task, error) {
		atomic.AddInt64(&worked, 1)
		time.Sleep(duration)
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

	if a := atomic.LoadInt64(&worked); a <= 295 {
		t.Errorf("more than 295-ish tasks should be worked, %d", a)
	}
}
