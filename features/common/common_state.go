package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

type CommonState interface{}

type commonState struct{}

func (c *commonState) ID() uint8 {
	return 5
}

func (c *commonState) Name() string {
	return "CommonState"
}

func (c *commonState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			1,
			"BatteryStateChanged",
			[]interface{}{
				uint8(0), // percent
			},
			c.batteryStateChanged,
		),
		arcommands.NewD2CCommand(
			7,
			"WifiSignalChanged",
			[]interface{}{
				int16(0), // rssi
			},
			c.wifiSignalChanged,
		),
	}
}

// TODO: Implement this
func (c *commonState) batteryStateChanged(args []interface{}) error {
	log.Debugf("the battery state changed: %v", args)
	return nil
}

// TODO: Implement this
func (c *commonState) wifiSignalChanged(args []interface{}) error {
	log.Debugf("the wifi signal changed: %v", args)
	return nil
}
