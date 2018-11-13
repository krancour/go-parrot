package common

import (
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/protocols/arnetwork"
)

// Feature ...
// TODO: Document this
type Feature interface {
	arcommands.D2CFeature
}

type feature struct {
	c2dChs map[uint8]chan<- arnetwork.Frame
}

// NewFeature ...
// TODO: Document this
func NewFeature(
	c2dChs map[uint8]chan<- arnetwork.Frame,
) Feature {
	return &feature{
		c2dChs: c2dChs,
	}
}

func (f *feature) ID() uint8 {
	return 0
}

func (f *feature) Name() string {
	return "common"
}

// TODO: Add stuff!
func (f *feature) D2CClasses() map[uint8]arcommands.D2CClass {
	return map[uint8]arcommands.D2CClass{

		5: arcommands.NewD2CClass(
			"CommonState",
			map[uint16]arcommands.D2CCommand{
				1: arcommands.NewD2CCommand(
					"BatteryStateChanged",
					[]interface{}{
						uint8(0), // percent
					},
					f.batteryStateChanged,
				),
				7: arcommands.NewD2CCommand(
					"WifiSignalChanged",
					[]interface{}{
						int16(0), // rssi
					},
					f.wifiSignalChanged,
				),
			},
		),
	}
}
