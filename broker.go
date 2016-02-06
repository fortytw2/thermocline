package thermocline

// NoVersion is used to denote non-versioned queues
var NoVersion = "n/a"

// A Broker is used as the communication interface between the workers and queue-er
type Broker interface {
	// Open returns a bidirectional channel through which tasks can be processed
	// or added.
	Open(queue string, version string) (chan *Task, error)

	// Failures returns how many failed tasks are on a given queue
	Failures(queue, version string) (int, error)
}
