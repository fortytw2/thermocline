package thermocline

// A Processor is the function given to any worker
type Processor func(*Task) ([]*Task, error)
