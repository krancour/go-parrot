package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Anti-flickering related states

const (
	ElectricFrequency50Hz int32 = 0
	ElectricFrequency60Hz int32 = 1

	AntiflickeringModeAuto int32 = 0
	AntiflickeringMode50Hz int32 = 1
	AntiflickeringMode60Hz int32 = 1
)

// AntiflickeringState ...
// TODO: Document this
type AntiflickeringState interface {
	lock.ReadLockable
	// ElectricFrequency returns the electrical frequency used by the camera. This
	// value maps to a constant and does not directly represent the frequency. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	ElectricFrequency() (int32, bool)
	// Mode returns the antiflickering mode. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	Mode() (int32, bool)
}

type antiflickeringState struct {
	sync.RWMutex
	electricFrequency *int32
	mode              *int32
}

func (a *antiflickeringState) ID() uint8 {
	return 30
}

func (a *antiflickeringState) Name() string {
	return "AntiflickeringState"
}

func (a *antiflickeringState) D2CCommands(
	log *log.Entry,
) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"electricFrequencyChanged",
			[]interface{}{
				int32(0), // frequency,
			},
			a.electricFrequencyChanged,
			log,
		),
		arcommands.NewD2CCommand(
			1,
			"modeChanged",
			[]interface{}{
				int32(0), // mode,
			},
			a.modeChanged,
			log,
		),
	}
}

// electricFrequencyChanged is invoked by the device when the electric frequency
// is changed.
func (a *antiflickeringState) electricFrequencyChanged(
	args []interface{},
	log *log.Entry,
) error {
	a.Lock()
	defer a.Unlock()
	a.electricFrequency = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"frequency", *a.electricFrequency,
	).Debug("electric frequency changed")
	return nil
}

// modeChanged is invoked by the device when the antiflickering mode is changed.
func (a *antiflickeringState) modeChanged(
	args []interface{},
	log *log.Entry,
) error {
	a.Lock()
	defer a.Unlock()
	a.mode = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"mode", *a.mode,
	).Debug("antiflickeringn mode changed")
	return nil
}

func (a *antiflickeringState) ElectricFrequency() (int32, bool) {
	if a.electricFrequency == nil {
		return 0, false
	}
	return *a.electricFrequency, true
}

func (a *antiflickeringState) Mode() (int32, bool) {
	if a.mode == nil {
		return 0, false
	}
	return *a.mode, true
}
