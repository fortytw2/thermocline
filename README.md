thermocline [![Build Status](https://travis-ci.org/fortytw2/thermocline.svg?branch=master)](https://travis-ci.org/fortytw2/thermocline) [![GoDoc](https://godoc.org/github.com/fortytw2/thermocline?status.svg)](http://godoc.org/github.com/fortytw2/thermocline) [![Go Report Card](https://goreportcard.com/badge/github.com/fortytw2/thermocline)](https://goreportcard.com/report/github.com/fortytw2/thermocline) [![codecov.io](https://codecov.io/github/fortytw2/thermocline/coverage.svg?branch=master)](https://codecov.io/github/fortytw2/thermocline?branch=master)
------

[DEPRECATED] - Using Channels in the core broker-interface was a poor decision in hindsight. ¯\_(ツ)_/¯ check out `github.com/fortytw2/hoplite`

A Library for implementing background-job-processing systems. Think of it as the implementation of the business-logic of `sidekiq`, without any convenience methods/helpful scheduling logic. Just raw workers, pools, and queues.


# Basic Usage

### Processor functions

A `Processor` is the function handed to a worker, of the type

```
type Processor func(*Task) ([]*Task, error)
```

Returning more tasks from a processor is optional, but useful in many cases. I would recommend implementing your `Processor` as a small wrapper that performs type-casting of `Task.Info`, an interface, to the type you actually want to work with.

### Pool Usage

You should be using `thermocline` through the `Pool` API, as using
individual workers may be a bit more error-prone/harder to work
with. Basic `Pool` use shown below.

```go
// A Broker is a simple message-queue implementation, an interface
// with only three functions to write to use your own.
var b thermocline.Broker
b = mem.NewBroker()

// create a new task every 10ms
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
// create a worker pool on the unversioned queue "test", with
p, err := thermocline.NewPool("test", thermocline.NoVersion, b, func(task *thermocline.Task) ([]*thermocline.Task, error) {
    atomic.AddInt64(&worked, 1)
    return nil, nil
}, 30)
if err != nil {
    t.Errorf("cannot create pool %s", err)
}

time.Sleep(500 * time.Millisecond)
ticker.Stop()
err = p.Stop()
if err != nil {
    // oh no a nasty error
    return err
}

fmt.Println(atomic.LoadInt64(&worked))
```

### Runtime tune-able pools

The number of workers in a given pool can be easily tuned at runtime
via `pool#Add(workers int)`, which can both add and remove workers
from a pool.

To stop a pool, simply call `pool.Stop()`, which will close out all workers and wait for their exit.


LICENSE
------

MIT
