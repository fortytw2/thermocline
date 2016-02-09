package thermocline_test

import (
	"testing"
	"time"

	"github.com/fortytw2/thermocline"
)

func TestDelayedTask(t *testing.T) {
	t.Parallel()

	_, reader, writer := SetupBroker(t)

	task, err := thermocline.NewDelayedTask("task infos", time.Second)
	if err != nil {
		t.Errorf("could not make delayed task, %s", err)
	}

	before := time.Now()
	writer <- task

	<-reader
	after := time.Now()

	if after.Sub(before) <= time.Second - time.Millisecond * 100 {
		t.Fatalf("task appeared before duration!")
	}
}

func TestLongDelayedTask(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping in short mode.")
	}

	_, reader, writer := SetupBroker(t)

	task, err := thermocline.NewDelayedTask("task infos", time.Second*10)
	if err != nil {
		t.Errorf("could not make delayed task, %s", err)
	}

	before := time.Now()
	writer <- task

	<-reader
	after := time.Now()

	if after.Sub(before) <= time.Second*9 {
		t.Fatalf("task appeared before duration!")
	}
}
