package thermocline

// NoVersion is used to denote non-versioned queues
var NoVersion = "n/a"

// A Broker is used as the communication interface between the workers and queue-er
type Broker interface {
	// Read returns a channel of incoming tasks
	Read(queue string, version string) (<-chan *Task, error)
	// Write returns a channel of outgoing tasks
	Write(queue string, version string) (chan<- *Task, error)
	// Stats returns current statistics of the broker (most likely just a bunch
	// of counters.)
	Stats() *Stats
}
