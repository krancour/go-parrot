package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

//

// FlightPlanSettingsState ...
// TODO: Document this
type FlightPlanSettingsState interface{}

type flightPlanSettingsState struct{}

func (f *flightPlanSettingsState) ID() uint8 {
	return 33
}

func (f *flightPlanSettingsState) Name() string {
	return "FlightPlanSettingsState"
}

func (f *flightPlanSettingsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"ReturnHomeOnDisconnectChanged",
			[]interface{}{
				uint8(0), // state,
				uint8(0), // isReadOnly,
			},
			f.returnHomeOnDisconnectChanged,
		),
	}
}

// TODO: Implement this
// Title: ReturnHome behavior during FlightPlan
// Description: Define behavior of drone when disconnection occurs during a
//   flight plan
// Support: 0901:4.1.0;090c:4.1.0;090e:1.4.0
// Triggered: by [setReturnHomeOnDisconnectMode](#0-32-0).
// Result:
func (f *flightPlanSettingsState) returnHomeOnDisconnectChanged(args []interface{}) error {
	// state := args[0].(uint8)
	//   1 if enabled, 0 if disabled
	// isReadOnly := args[1].(uint8)
	//   1 if readOnly, 0 if writable
	log.Info("common.returnHomeOnDisconnectChanged() called")
	return nil
}
