package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Network state from Product

// NetworkState ...
// TODO: Document this
type NetworkState interface {
	lock.ReadLockable
}

type networkState struct {
	sync.RWMutex
}

func (n *networkState) ID() uint8 {
	return 14
}

func (n *networkState) Name() string {
	return "NetworkState"
}

func (n *networkState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"WifiScanListChanged",
			[]interface{}{
				string(0), // ssid,
				int16(0),  // rssi,
				int32(0),  // band,
				uint8(0),  // channel,
			},
			n.wifiScanListChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"AllWifiScanChanged",
			[]interface{}{},
			n.allWifiScanChanged,
		),
		arcommands.NewD2CCommand(
			2,
			"WifiAuthChannelListChanged",
			[]interface{}{
				int32(0), // band,
				uint8(0), // channel,
				uint8(0), // in_or_out,
			},
			n.wifiAuthChannelListChanged,
		),
		arcommands.NewD2CCommand(
			3,
			"AllWifiAuthChannelChanged",
			[]interface{}{},
			n.allWifiAuthChannelChanged,
		),
	}
}

// TODO: Implement this
// Title: Wifi scan results
// Description: Wifi scan results.\n Please note that the list is not complete
//   until you receive the event [WifiScanEnded](#1-14-1).
// Support: 0901;090c;090e
// Triggered: for each wifi network scanned after a [ScanWifi](#1-13-0)
// Result:
func (n *networkState) wifiScanListChanged(args []interface{}) error {
	// ssid := args[0].(string)
	//   SSID of the AP
	// rssi := args[1].(int16)
	//   RSSI of the AP in dbm (negative value)
	// band := args[2].(int32)
	//   The band : 2.4 GHz or 5 GHz
	//   0: 2_4ghz: 2.4 GHz band
	//   1: 5ghz: 5 GHz band
	// channel := args[3].(uint8)
	//   Channel of the AP
	log.Info("ardrone3.wifiScanListChanged() called")
	return nil
}

// TODO: Implement this
// Title: Wifi scan ended
// Description: Wifi scan ended.\n When receiving this event, the list of
//   [WifiScanResults](#1-14-0) is complete.
// Support: 0901;090c;090e
// Triggered: after the last [WifiScanResult](#1-14-0) has been sent.
// Result:
func (n *networkState) allWifiScanChanged(args []interface{}) error {
	log.Info("ardrone3.allWifiScanChanged() called")
	return nil
}

// TODO: Implement this
// Title: Available wifi channels
// Description: Available wifi channels.\n Please note that the list is not
//   complete until you receive the event
//   [AvailableWifiChannelsCompleted](#1-14-3).
// Support: 0901;090c;090e
// Triggered: for each available channel after a
//   [GetAvailableWifiChannels](#1-13-1).
// Result:
func (n *networkState) wifiAuthChannelListChanged(args []interface{}) error {
	// band := args[0].(int32)
	//   The band of this channel : 2.4 GHz or 5 GHz
	//   0: 2_4ghz: 2.4 GHz band
	//   1: 5ghz: 5 GHz band
	// channel := args[1].(uint8)
	//   The authorized channel.
	// in_or_out := args[2].(uint8)
	//   Bit 0 is 1 if channel is authorized outside (0 otherwise) ; Bit 1 is 1 if
	//   channel is authorized inside (0 otherwise)
	log.Info("ardrone3.wifiAuthChannelListChanged() called")
	return nil
}

// TODO: Implement this
// Title: Available wifi channels completed
// Description: Available wifi channels completed.\n When receiving this event,
//   the list of [AvailableWifiChannels](#1-14-2) is complete.
// Support: 0901;090c;090e
// Triggered: after the last [AvailableWifiChannel](#1-14-2) has been sent.
// Result:
func (n *networkState) allWifiAuthChannelChanged(args []interface{}) error {
	log.Info("ardrone3.allWifiAuthChannelChanged() called")
	return nil
}
