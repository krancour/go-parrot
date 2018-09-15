package bebop2

import (
	"fmt"

	"github.com/krancour/go-parrot/features/ardrone3"
	"github.com/krancour/go-parrot/features/common"
	"github.com/krancour/go-parrot/protocols/arnetwork"
	"github.com/krancour/go-parrot/protocols/arnetworkal/wifi"
)

// Client ...
// TODO: Document this
type Client interface {
	// TODO: Add all supported commands to this interface
	// The first ones we experiment with should be configuration related and
	// not anything to do with flight. We want to be pretty sure that everything
	// works before we try flying!
}

type client struct {
	common   common.Feature
	ardrone3 ardrone3.Feature
}

// NewClient ...
// TODO: Document this
func NewClient() (Client, error) {
	conn, err := wifi.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("error creating wifi connection: %s", err)
	}
	bufMan, err := arnetwork.NewBufferManager(
		conn,
		// TODO: Add device-specific buffers here
	)
	if err != nil {
		return nil, fmt.Errorf("error creating buffer manager: %s", err)
	}
	return &client{
		common:   common.NewFeature(bufMan),
		ardrone3: ardrone3.NewFeature(bufMan),
	}, nil
}
