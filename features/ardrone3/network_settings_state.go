package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Network settings state from product

const (
	WifiTypeAutoAll        int32 = 0
	WifiTypeAuto2Point4GHz int32 = 1
	WifiTypeAuto5GHz       int32 = 2
	WifiTypeManual         int32 = 3

	WifiBand2Point4GHz int32 = 0
	WifiBand5GHz       int32 = 1
	WifiBandAll        int32 = 2

	WifiSecurityTypeOpen int32 = 0
	WifiSecurityTypeWPA2 int32 = 1

	WifiSecurityKeyTypePlaintext int32 = 0
)

// NetworkSettingsState ...
// TODO: Document this
type NetworkSettingsState interface {
	lock.ReadLockable
	// Type returns the type of wifi selection settings. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	Type() (int32, bool)
	// Band returns the actual wifi band state. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	Band() (int32, bool)
	// Channel returns the wifi channel (depends on the band). A boolean value is
	// also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	Channel() (uint8, bool)
	// SecurityType returns the type of wifi security in use. A boolean value is
	// also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	SecurityType() (int32, bool)
	// SecurityKey returns the wifi security key. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	SecurityKey() (string, bool)
	// SecurityKeyType returns the type of wifi security key in use. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	SecurityKeyType() (int32, bool)
}

type networkSettingsState struct {
	sync.RWMutex
	// tipe represents the type of wifi selection settings
	tipe *int32
	// band represents actual wifi band state
	band *int32
	// channel is the wifi channel (depends of the band)
	channel *uint8
	// securityType represents the type of wifi security in use
	securityType *int32
	// securityKey is the wifi security key
	securityKey *string
	// securityKeyType represents the type of wifi security key in use
	securityKeyType *int32
}

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
		arcommands.NewD2CCommand(
			1,
			"wifiSecurityChanged",
			[]interface{}{
				int32(0), // type,
			},
			n.wifiSecurityChanged,
		),
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

// wifiSelectionChanged is invoked by the device when wifi selection changes.
func (n *networkSettingsState) wifiSelectionChanged(args []interface{}) error {
	n.Lock()
	defer n.Unlock()
	n.tipe = ptr.ToInt32(args[0].(int32))
	n.band = ptr.ToInt32(args[1].(int32))
	n.channel = ptr.ToUint8(args[2].(uint8))
	log.WithField(
		"type", *n.tipe,
	).WithField(
		"band", *n.band,
	).WithField(
		"channel", *n.channel,
	).Debug("wifi selection changed")
	return nil
}

// wifiSecurityChanged is deprecated in favor of wifiSecurity, but since we can
// still see this command being invoked, we'll implement the command to avoid a
// warning, but the implementation will remain a no-op unless / until such time
// that it becomes clear that older versions of the firmware might require us to
// support both commands.
func (n *networkSettingsState) wifiSecurityChanged(args []interface{}) error {
	// type := args[0].(int32)
	//   The type of wifi security (open, wpa2)
	//   0: open: Wifi is not protected by any security (default)
	//   1: wpa2: Wifi is protected by wpa2
	log.Debug("wifi security changed-- this is a no-op")
	return nil
}

// wifiSecurity is invoked by the device when wifi security changes.
func (n *networkSettingsState) wifiSecurity(args []interface{}) error {
	n.Lock()
	defer n.Unlock()
	n.securityType = ptr.ToInt32(args[0].(int32))
	n.securityKey = ptr.ToString(args[1].(string))
	n.securityKeyType = ptr.ToInt32(args[2].(int32))
	log.WithField(
		"type", *n.securityType,
	).WithField(
		"key", *n.securityKey,
	).WithField(
		"keyType", *n.securityKeyType,
	).Debug("wifi security changed")
	return nil
}

func (n *networkSettingsState) Type() (int32, bool) {
	if n.tipe == nil {
		return 0, false
	}
	return *n.tipe, true
}

func (n *networkSettingsState) Band() (int32, bool) {
	if n.band == nil {
		return 0, false
	}
	return *n.band, true
}

func (n *networkSettingsState) Channel() (uint8, bool) {
	if n.channel == nil {
		return 0, false
	}
	return *n.channel, true
}

func (n *networkSettingsState) SecurityType() (int32, bool) {
	if n.securityType == nil {
		return 0, false
	}
	return *n.securityType, true
}

func (n *networkSettingsState) SecurityKey() (string, bool) {
	if n.securityKey == nil {
		return "", false
	}
	return *n.securityKey, true
}

func (n *networkSettingsState) SecurityKeyType() (int32, bool) {
	if n.securityKeyType == nil {
		return 0, false
	}
	return *n.securityKeyType, true
}
