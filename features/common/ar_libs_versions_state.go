package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// ARlibs Versions Commands

// ARLibsVersionsState ...
// TODO: Document this
type ARLibsVersionsState interface{}

type arLibsVersionsState struct{}

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

// TODO: Implement this
func (a *arLibsVersionsState) deviceLibARCommandsVersion(
	args []interface{},
) error {
	// version := args[0].(string)
	//   version of libARCommands (&#34;1.2.3.4&#34; format)
	log.Info("common.deviceLibARCommandsVersion() called")
	return nil
}
