package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// FlightPlan Event commands

// FlightPlanEvent ...
// TODO: Document this
type FlightPlanEvent interface {
	lock.ReadLockable
}

type flightPlanEvent struct {
	sync.RWMutex
}

func newFlightPlanEvent() *flightPlanEvent {
	return &flightPlanEvent{}
}

func (f *flightPlanEvent) ClassID() uint8 {
	return 19
}

func (f *flightPlanEvent) ClassName() string {
	return "FlightPlanEvent"
}

func (f *flightPlanEvent) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"StartingErrorEvent",
			[]interface{}{},
			f.startingErrorEvent,
			log,
		),
		// arcommands.NewD2CCommand(
		// 	1,
		// 	"SpeedBridleEvent",
		// 	[]interface{}{},
		// 	f.speedBridleEvent,
		// ),
	}
}

// TODO: Implement this
// Title: FlightPlan start error
// Description: FlightPlan start error.\n\n **This event is a notification, you
//   can&#39;t retrieve it in the cache of the device controller.**
// Support: 0901:2.0.29;090c;090e
// Triggered: on an error after a [StartFlightPlan](#0-11-0).
// Result:
func (f *flightPlanEvent) startingErrorEvent(
	args []interface{},
	log *log.Entry,
) error {
	log.Warn("command not implemented")
	return nil
}

// // TODO: Implement this
// // Title: FlightPlan speed clamping
// // Description: FlightPlan speed clamping.\n Sent when a speed specified in the
// //   FlightPlan file is considered too high by the drone.\n\n **This event is a
// //   notification, you can&#39;t retrieve it in the cache of the device
// //   controller.**
// // Support: none
// // Triggered: on an speed related clamping after a [StartFlightPlan](#0-11-0).
// // Result:
// func (f *flightPlanEvent) speedBridleEvent(args []interface{}) error {
// 	log.Warn("command not implemented")
// 	return nil
// }
