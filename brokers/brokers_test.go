package brokers_test

import (
	"testing"

	"github.com/fortytw2/thermocline"
	"github.com/fortytw2/thermocline/brokers/mem"
)

var brokers = []thermocline.Broker{mem.NewBroker()}

func TestBrokers(t *testing.T) {

}
