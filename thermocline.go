package thermocline

// Work automates creation and deletion of workers
func Work(pf Processor, workers int, stop chan struct{}) error {

	return nil
}

type StartReq struct {
	Functions map[string]map[string]Processor
	Workers   int
}
