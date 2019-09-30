package registry

import (
	"fmt"
	"sync"
)

// Registry is the local service registry.
// Here the MCP maintains the information of the service endpoints
// on which to perform health check operations.
type Registry struct{}

var instance *Registry
var onceRegistry sync.Once

// GetRegistry returns the Registry instance. The Registry instance is
// initialized only one time, as a singleton.
func GetRegistry() *Registry {
	onceRegistry.Do(func() {
		instance = &Registry{}
	})

	return instance
}

// GetServiceInfo returns the service information for the service identified by the <name> in input.
// Returns a nil if no <name> service is registered with the local registry.
func (r *Registry) GetServiceInfo(name string) string {
	return fmt.Sprintf("%s service info", name)
}
