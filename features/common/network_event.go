package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Network Event from product

// NetworkEvent ...
// TODO: Document this
type NetworkEvent interface {
	lock.ReadLockable
}

type networkEvent struct {
	sync.RWMutex
}

func (n *networkEvent) ID() uint8 {
	return 1
}

func (n *networkEvent) Name() string {
	return "NetworkEvent"
}

func (n *networkEvent) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"Disconnection",
			[]interface{}{
				int32(0), // cause,
			},
			n.disconnection,
			log,
		),
	}
}

// TODO: Implement this
// Title: Drone will disconnect
// Description: Drone will disconnect.\n This event is mainly triggered when the
//   user presses on the power button of the product.\n\n **This event is a
//   notification, you can&#39;t retrieve it in the cache of the device
//   controller.**
// Support: 0901;090c
// Triggered: mainly when the user presses the power button of the drone.
// Result:
func (n *networkEvent) disconnection(args []interface{}, log *log.Entry) error {
	// cause := args[0].(int32)
	//   Cause of the disconnection of the product
	//   0: off_button: The button off has been pressed
	//   1: unknown: Unknown generic cause
	log.Warn("command not implemented")
	return nil
}
