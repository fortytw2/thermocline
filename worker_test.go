package thermocline_test

import (
	"errors"
	"log"
	"testing"

	"github.com/fortytw2/thermocline"
	"github.com/fortytw2/thermocline/brokers/mem"
)

func TestWorker(t *testing.T) {
	broker := mem.NewBroker()

	tasks, err := broker.Open("test", thermocline.NoVersion)
	if err != nil {
		t.Errorf("could not open queue '%s'", err)
	}

	task, err := thermocline.NewTask("test")
	if err != nil {
		t.Error("could not create test task", err)
	}

	stopper := make(chan struct{})
	w := thermocline.NewWorker(tasks, func(task *thermocline.Task) ([]*thermocline.Task, error) {
		log.Printf("%+v", task)
		return nil, nil
	}, stopper)

	go w.Work()
	tasks <- task

	close(stopper)
}

func TestWorkerRetries(t *testing.T) {
	broker := mem.NewBroker()

	tasks, err := broker.Open("test", thermocline.NoVersion)
	if err != nil {
		t.Errorf("could not open queue '%s'", err)
	}

	task, err := thermocline.NewTask("test")
	if err != nil {
		t.Error("could not create test task", err)
	}

	stopper := make(chan struct{})
	w := thermocline.NewWorker(tasks, func(task *thermocline.Task) ([]*thermocline.Task, error) {
		return nil, errors.New("cannot proccess task")
	}, stopper)

	go w.Work()
	tasks <- task

	close(stopper)
}

func TestWorkerStop(t *testing.T) {

}
