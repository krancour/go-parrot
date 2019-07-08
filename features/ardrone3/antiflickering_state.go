package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Anti-flickering related states

// AntiflickeringState ...
// TODO: Document this
type AntiflickeringState interface {
	lock.ReadLockable
}

type antiflickeringState struct {
	sync.RWMutex
}

func (a *antiflickeringState) ID() uint8 {
	return 30
}

func (a *antiflickeringState) Name() string {
	return "AntiflickeringState"
}

func (a *antiflickeringState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"electricFrequencyChanged",
			[]interface{}{
				int32(0), // frequency,
			},
			a.electricFrequencyChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"modeChanged",
			[]interface{}{
				int32(0), // mode,
			},
			a.modeChanged,
		),
	}
}

// TODO: Implement this
// Title: Electric frequency
// Description: Electric frequency.\n This piece of information is used for the
//   antiflickering when the [AntiflickeringMode](#1-30-1) is set to *auto*.
// Support: 0901;090c
// Triggered: by [SetElectricFrequency](#1-29-0).
// Result:
func (a *antiflickeringState) electricFrequencyChanged(
	args []interface{},
) error {
	// frequency := args[0].(int32)
	//   Type of the electric frequency
	//   0: fiftyHertz: Electric frequency of the country is 50hz
	//   1: sixtyHertz: Electric frequency of the country is 60hz
	log.Info("ardrone3.electricFrequencyChanged() called")
	return nil
}

// TODO: Implement this
// Title: Antiflickering mode
// Description: Antiflickering mode.
// Support: 0901;090c
// Triggered: by [SetAntiflickeringMode](#1-29-1).
// Result:
func (a *antiflickeringState) modeChanged(args []interface{}) error {
	// mode := args[0].(int32)
	//   Mode of the anti flickering functionnality
	//   0: auto: Anti flickering based on the electric frequency previously sent
	//   1: FixedFiftyHertz: Anti flickering based on a fixed frequency of 50Hz
	//   2: FixedSixtyHertz: Anti flickering based on a fixed frequency of 60Hz
	log.Info("ardrone3.modeChanged() called")
	return nil
}
