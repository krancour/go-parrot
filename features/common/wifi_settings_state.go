package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Wifi settings state from product

// WifiSettingsState ...
// TODO: Document this
type WifiSettingsState interface{}

type wifiSettingsState struct{}

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

// TODO: Implement this
// Title: Wifi outdoor mode
// Description: Wifi outdoor mode.
// Support: 0901;0902;0905;0906;090c;090e
// Triggered: by [SetWifiOutdoorMode](#0-9-0).
// Result:
func (w *wifiSettingsState) outdoorSettingsChanged(args []interface{}) error {
	// outdoor := args[0].(uint8)
	//   1 if it should use outdoor wifi settings, 0 otherwise
	log.Info("common.outdoorSettingsChanged() called")
	return nil
}
