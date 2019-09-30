package discoverer

import (
	"log"
	"time"

	"github.com/itross/sgulmcp/internal/registry"
)

// Discoverer is the service discovery agent.
// It calls the SgulREG service registry and get the full registry.
type Discoverer struct {
	registry *registry.Registry
}

// New returns a new Discoverer instance.
func New() *Discoverer {
	return &Discoverer{
		registry: registry.GetRegistry(),
	}
}

// Discover .
func (d *Discoverer) Discover() error {
	time.Sleep(5 * time.Second)
	log.Printf("%s discovered", d.registry.GetServiceInfo("test-service"))
	// return errors.New("FAKE DISCOVERY ERROR")
	return nil
}
