package ardrone3

import "github.com/krancour/go-parrot/protocols/arnetwork"

// Feature ...
// TODO: Document this
type Feature interface{}

type feature struct {
	c2dChs map[uint8]chan<- arnetwork.Frame
	d2cChs map[uint8]<-chan arnetwork.Frame
}

// NewFeature ...
// TODO: Document this
func NewFeature(
	c2dChs map[uint8]chan<- arnetwork.Frame,
	d2cChs map[uint8]<-chan arnetwork.Frame,
) Feature {
	return &feature{
		c2dChs: c2dChs,
		d2cChs: d2cChs,
	}
}
