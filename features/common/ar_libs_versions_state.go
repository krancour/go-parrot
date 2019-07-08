package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// ARlibs Versions Commands

// ARLibsVersionsState ...
// TODO: Document this
type ARLibsVersionsState interface {
	lock.ReadLockable
	// DeviceARLibsVersion returns the version of the ARLibs library in use by the
	// devcie. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	DeviceARLibsVersion() (string, bool)
}

type arLibsVersionsState struct {
	sync.RWMutex
	deviceARLibsVersion *string
}

func (a *arLibsVersionsState) ID() uint8 {
	return 18
}

func (a *arLibsVersionsState) Name() string {
	return "ARLibsVersionsState"
}

func (a *arLibsVersionsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"ControllerLibARCommandsVersion",
			[]interface{}{
				string(0), // version,
			},
			a.controllerLibARCommandsVersion,
		),
		arcommands.NewD2CCommand(
			1,
			"SkyControllerLibARCommandsVersion",
			[]interface{}{
				string(0), // version,
			},
			a.skyControllerLibARCommandsVersion,
		),
		arcommands.NewD2CCommand(
			2,
			"DeviceLibARCommandsVersion",
			[]interface{}{
				string(0), // version,
			},
			a.deviceLibARCommandsVersion,
		),
	}
}

// TODO: Implement this
func (a *arLibsVersionsState) controllerLibARCommandsVersion(
	args []interface{},
) error {
	// version := args[0].(string)
	//   version of libARCommands (&#34;1.2.3.4&#34; format)
	log.Info("common.controllerLibARCommandsVersion() called")
	return nil
}

// TODO: Implement this
func (a *arLibsVersionsState) skyControllerLibARCommandsVersion(
	args []interface{},
) error {
	// version := args[0].(string)
	//   version of libARCommands (&#34;1.2.3.4&#34; format)
	log.Info("common.skyControllerLibARCommandsVersion() called")
	return nil
}

// deviceLibARCommandsVersion is invoked by the device to indicate what version
// of the ARLibs library it is using.
func (a *arLibsVersionsState) deviceLibARCommandsVersion(
	args []interface{},
) error {
	a.Lock()
	defer a.Unlock()
	a.deviceARLibsVersion = ptr.ToString(args[0].(string))
	log.WithField(
		"deviceARLibsVersion", *a.deviceARLibsVersion,
	).Debug("device ARLibs version changed")
	return nil
}

func (a *arLibsVersionsState) DeviceARLibsVersion() (string, bool) {
	if a.deviceARLibsVersion == nil {
		return "", false
	}
	return *a.deviceARLibsVersion, true
}
