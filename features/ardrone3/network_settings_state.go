package ardrone3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Network settings state from product

// NetworkSettingsState ...
// TODO: Document this
type NetworkSettingsState interface{}

type networkSettingsState struct{}

func (n *networkSettingsState) ID() uint8 {
	return 10
}

func (n *networkSettingsState) Name() string {
	return "NetworkSettingsState"
}

func (n *networkSettingsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"WifiSelectionChanged",
			[]interface{}{
				int32(0), // type,
				int32(0), // band,
				uint8(0), // channel,
			},
			n.wifiSelectionChanged,
		),
		// arcommands.NewD2CCommand(
		// 	1,
		// 	"wifiSecurityChanged",
		// 	[]interface{}{
		// 		int32(0), // type,
		// 	},
		// 	n.wifiSecurityChanged,
		// ),
		arcommands.NewD2CCommand(
			2,
			"wifiSecurity",
			[]interface{}{
				int32(0),  // type,
				string(0), // key,
				int32(0),  // keyType,
			},
			n.wifiSecurity,
		),
	}
}

// TODO: Implement this
// Title: Wifi selection
// Description: Wifi selection.
// Support: 0901;090c;090e
// Triggered: by [SelectWifi](#1-9-0).
// Result:
func (n *networkSettingsState) wifiSelectionChanged(args []interface{}) error {
	// type := args[0].(int32)
	//   The type of wifi selection settings
	//   0: auto_all: Auto selection
	//   1: auto_2_4ghz: Auto selection 2.4ghz
	//   2: auto_5ghz: Auto selection 5 ghz
	//   3: manual: Manual selection
	// band := args[1].(int32)
	//   The actual wifi band state
	//   0: 2_4ghz: 2.4 GHz band
	//   1: 5ghz: 5 GHz band
	//   2: all: Both 2.4 and 5 GHz bands
	// channel := args[2].(uint8)
	//   The channel (depends of the band)
	log.Info("ardrone3.wifiSelectionChanged() called")
	return nil
}

// // TODO: Implement this
// // Title: Wifi security type
// // Description: Wifi security type.
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (n *networkSettingsState) wifiSecurityChanged(args []interface{}) error {
// 	// type := args[0].(int32)
// 	//   The type of wifi security (open, wpa2)
// 	//   0: open: Wifi is not protected by any security (default)
// 	//   1: wpa2: Wifi is protected by wpa2
// 	log.Info("ardrone3.wifiSecurityChanged() called")
// 	return nil
// }

// TODO: Implement this
// Title: Wifi security type
// Description: Wifi security type.
// Support: 0901;090c;090e
// Triggered: by [SetWifiSecurityType](#1-9-1).
// Result:
func (n *networkSettingsState) wifiSecurity(args []interface{}) error {
	// type := args[0].(int32)
	//   The type of wifi security (open, wpa2)
	//   0: open: Wifi is not protected by any security (default)
	//   1: wpa2: Wifi is protected by wpa2
	// key := args[1].(string)
	//   The key used to secure the network (empty if type is open)
	// keyType := args[2].(int32)
	//   Type of the key
	//   0: plain: Key is plain text, not encrypted
	log.Info("ardrone3.wifiSecurity() called")
	return nil
}
