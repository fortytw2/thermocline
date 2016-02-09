package thermocline_test

import (
	"testing"
	"time"

	"github.com/fortytw2/thermocline"
	"github.com/fortytw2/thermocline/brokers/mem"
)

func TestDelayedTask(t *testing.T) {
	t.Parallel()

	var b thermocline.Broker
	b = mem.NewBroker()

	reader, err := b.Read("test", thermocline.NoVersion)
	if err != nil {
		t.Errorf("could not open queue '%s'", err)
	}

	writer, err := b.Write("test", thermocline.NoVersion)
	if err != nil {
		t.Errorf("could not open queue '%s'", err)
	}

	task, err := thermocline.NewDelayedTask("task infos", time.Second)
	if err != nil {
		t.Errorf("could not make delayed task, %s", err)
	}

	before := time.Now()
	writer <- task

	<-reader
	after := time.Now()

	if after.Sub(before) <= time.Second {
		t.Fatalf("task appeared before duration!")
	}
}

func TestLongDelayedTask(t *testing.T) {
	t.Parallel()

	var b thermocline.Broker
	b = mem.NewBroker()

	reader, err := b.Read("test", thermocline.NoVersion)
	if err != nil {
		t.Errorf("could not open queue '%s'", err)
	}

	writer, err := b.Write("test", thermocline.NoVersion)
	if err != nil {
		t.Errorf("could not open queue '%s'", err)
	}

	task, err := thermocline.NewDelayedTask("task infos", time.Second*10)
	if err != nil {
		t.Errorf("could not make delayed task, %s", err)
	}

	before := time.Now()
	writer <- task

	<-reader
	after := time.Now()

	if after.Sub(before) <= time.Second*10 {
		t.Fatalf("task appeared before duration!")
	}
}
