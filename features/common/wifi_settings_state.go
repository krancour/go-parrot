package common

import (
	"sync"

	"github.com/krancour/go-parrot/lock"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Wifi settings state from product

// WifiSettingsState ...
// TODO: Document this
type WifiSettingsState interface {
	lock.ReadLockable
	// Outdoors returns a boolean indicating whether the device is using outdoor
	// wifi settings. A boolean value is also returned, indicating whether the
	// first value was reported by the device (true) or a default value (false).
	// This permits callers to distinguish real zero values from default zero
	// values.
	Outdoors() (bool, bool)
}

type wifiSettingsState struct {
	sync.RWMutex
	outdoors *bool
}

func newWifiSettingsState() *wifiSettingsState {
	return &wifiSettingsState{}
}

func (w *wifiSettingsState) ID() uint8 {
	return 10
}

func (w *wifiSettingsState) Name() string {
	return "WifiSettingsState"
}

func (w *wifiSettingsState) D2CCommands(
	log *log.Entry,
) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"outdoorSettingsChanged",
			[]interface{}{
				uint8(0), // outdoor,
			},
			w.outdoorSettingsChanged,
			log,
		),
	}
}

// Invoked by the device to indicate whether it is using outdoor wifi settings.
func (w *wifiSettingsState) outdoorSettingsChanged(
	args []interface{},
	log *log.Entry,
) error {
	w.Lock()
	defer w.Unlock()
	w.outdoors = ptr.ToBool(args[0].(uint8) == 1)
	log.WithField(
		"outdoors", *w.outdoors,
	).Debug("outdoor settings changed")
	return nil
}

func (w *wifiSettingsState) Outdoors() (bool, bool) {
	if w.outdoors == nil {
		return false, false
	}
	return *w.outdoors, true
}
