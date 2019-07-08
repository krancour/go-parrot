package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Events of Piloting

// PilotingEvent ...
// TODO: Document this
type PilotingEvent interface {
	lock.ReadLockable
}

type pilotingEvent struct {
	sync.RWMutex
}

func (p *pilotingEvent) ID() uint8 {
	return 34
}

func (p *pilotingEvent) Name() string {
	return "PilotingEvent"
}

func (p *pilotingEvent) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"moveByEnd",
			[]interface{}{
				float32(0), // dX,
				float32(0), // dY,
				float32(0), // dZ,
				float32(0), // dPsi,
				int32(0),   // error,
			},
			p.moveByEnd,
		),
	}
}

// TODO: Implement this
// Title: Relative move ended
// Description: Relative move ended.\n Informs about the move that the drone
//   managed to do and why it stopped.
// Support: 0901:3.3.0;090c:3.3.0
// Triggered: when the drone reaches its target or when it is interrupted by
//   another [moveBy command](#1-0-7) or when an error occurs.
// Result:
func (p *pilotingEvent) moveByEnd(args []interface{}) error {
	// dX := args[0].(float32)
	//   Distance traveled along the front axis [m]
	// dY := args[1].(float32)
	//   Distance traveled along the right axis [m]
	// dZ := args[2].(float32)
	//   Distance traveled along the down axis [m]
	// dPsi := args[3].(float32)
	//   Applied angle on heading [rad]
	// error := args[4].(int32)
	//   Error to explain the event
	//   0: ok: No Error ; The relative displacement
	//   1: unknown: Unknown generic error
	//   2: busy: The Device is busy ; command moveBy ignored
	//   3: notAvailable: Command moveBy is not available ; command moveBy ignored
	//   4: interrupted: Command moveBy interrupted
	log.Info("ardrone3.oveByEnd() called")
	return nil
}
