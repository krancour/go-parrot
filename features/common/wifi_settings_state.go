package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Wifi settings state from product

// WifiSettingsState ...
// TODO: Document this
type WifiSettingsState interface {
	// RLock blocks until a read lock is obtained. This permits callers to procede
	// with querying any or all attributes of the wifi settings state without
	// worry that some attributes will be overwritten as others are read. i.e. It
	// permits the possibility of taking an atomic snapshop of wifi settings
	// state. Note that use of this function is not obligatory for applications
	// that do not require such guarantees. Callers MUST call RUnlock() or else
	// wifi settings state will never resume updating.
	RLock()
	// RUnlock releases a read lock on the wifi settings state. See RLock().
	RUnlock()
	// Outdoors returns a boolean indicating whether the device is using outdoor
	// wifi settings. A boolean value is also returned, indicating whether the
	// first value was reported by the device (true) or a default value (false).
	// This permits callers to distinguish real zero values from default zero
	// values.
	Outdoors() (bool, bool)
}

type wifiSettingsState struct {
	outdoors *bool
	lock     sync.RWMutex
}

func (w *wifiSettingsState) ID() uint8 {
	return 10
}

func (w *wifiSettingsState) Name() string {
	return "WifiSettingsState"
}

func (w *wifiSettingsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"outdoorSettingsChanged",
			[]interface{}{
				uint8(0), // outdoor,
			},
			w.outdoorSettingsChanged,
		),
	}
}

// Invoked by the device to indicate whether it is using outdoor wifi settings.
func (w *wifiSettingsState) outdoorSettingsChanged(args []interface{}) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.outdoors = ptr.ToBool(args[0].(uint8) == 1)
	log.WithField(
		"outdoors", *w.outdoors,
	).Debug("outdoor settings changed")
	return nil
}

func (w *wifiSettingsState) RLock() {
	w.lock.RLock()
}

func (w *wifiSettingsState) RUnlock() {
	w.lock.RUnlock()
}

func (w *wifiSettingsState) Outdoors() (bool, bool) {
	if w.outdoors == nil {
		return false, false
	}
	return *w.outdoors, true
}
