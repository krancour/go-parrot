package ardrone3

import "github.com/krancour/drone-examples/pkg/parrot/protocols/arnetwork"

// Feature ...
// TODO: Document this
type Feature interface{}

type feature struct {
	bufMan arnetwork.BufferManager
}

// NewFeature ...
// TODO: Document this
func NewFeature(bufferManager arnetwork.BufferManager) Feature {
	return &feature{
		bufMan: bufferManager,
	}
}
