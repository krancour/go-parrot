package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Settings state from product

// SettingsState ...
// TODO: Document this
type SettingsState interface {
	lock.ReadLockable
}

type settingsState struct {
	lock sync.RWMutex
}

func (s *settingsState) ID() uint8 {
	return 16
}

func (s *settingsState) Name() string {
	return "SettingsState"
}

func (s *settingsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		// arcommands.NewD2CCommand(
		// 	0,
		// 	"ProductMotorVersionListChanged",
		// 	[]interface{}{
		// 		uint8(0),  // motor_number,
		// 		string(0), // type,
		// 		string(0), // software,
		// 		string(0), // hardware,
		// 	},
		// 	s.productMotorVersionListChanged,
		// ),
		arcommands.NewD2CCommand(
			1,
			"ProductGPSVersionChanged",
			[]interface{}{
				string(0), // software,
				string(0), // hardware,
			},
			s.productGPSVersionChanged,
		),
		arcommands.NewD2CCommand(
			2,
			"MotorErrorStateChanged",
			[]interface{}{
				uint8(0), // motorIds,
				int32(0), // motorError,
			},
			s.motorErrorStateChanged,
		),
		arcommands.NewD2CCommand(
			3,
			"MotorSoftwareVersionChanged",
			[]interface{}{
				string(0), // version,
			},
			s.motorSoftwareVersionChanged,
		),
		arcommands.NewD2CCommand(
			4,
			"MotorFlightsStatusChanged",
			[]interface{}{
				uint16(0), // nbFlights,
				uint16(0), // lastFlightDuration,
				uint32(0), // totalFlightDuration,
			},
			s.motorFlightsStatusChanged,
		),
		arcommands.NewD2CCommand(
			5,
			"MotorErrorLastErrorChanged",
			[]interface{}{
				int32(0), // motorError,
			},
			s.motorErrorLastErrorChanged,
		),
		arcommands.NewD2CCommand(
			6,
			"P7ID",
			[]interface{}{
				string(0), // serialID,
			},
			s.p7ID,
		),
		arcommands.NewD2CCommand(
			7,
			"CPUID",
			[]interface{}{
				string(0), // id,
			},
			s.cPUID,
		),
	}
}

// // TODO: Implement this
// // Title: Motor version
// // Description: Motor version.
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (s *settingsState) productMotorVersionListChanged(
// 	args []interface{},
// ) error {
// 	// motor_number := args[0].(uint8)
// 	//   Product Motor number
// 	// type := args[1].(string)
// 	//   Product Motor type
// 	// software := args[2].(string)
// 	//   Product Motors software version
// 	// hardware := args[3].(string)
// 	//   Product Motors hardware version
// 	log.Info("ardrone3.productMotorVersionListChanged() called")
// 	return nil
// }

// TODO: Implement this
// Title: GPS version
// Description: GPS version.
// Support: 0901;090c;090e
// Triggered: at connection.
// Result:
func (s *settingsState) productGPSVersionChanged(args []interface{}) error {
	// software := args[0].(string)
	//   Product GPS software version
	// hardware := args[1].(string)
	//   Product GPS hardware version
	log.Info("ardrone3.productGPSVersionChanged() called")
	return nil
}

// TODO: Implement this
// Title: Motor error
// Description: Motor error.\n This event is sent back to *noError* as soon as
//   the motor error disappear. To get the last motor error, see
//   [LastMotorError](#1-16-5)
// Support: 0901;090c;090e
// Triggered: when a motor error occurs.
// Result:
func (s *settingsState) motorErrorStateChanged(args []interface{}) error {
	// motorIds := args[0].(uint8)
	//   Bit field for concerned motor. If bit 0 = 1, motor 1 is affected by this
	//   error. Same with bit 1, 2 and 3. Motor 1: front left Motor 2: front right
	//   Motor 3: back right Motor 4: back left
	// motorError := args[1].(int32)
	//   Enumeration of the motor error
	//   0: noError: No error detected
	//   1: errorEEPRom: EEPROM access failure
	//   2: errorMotorStalled: Motor stalled
	//   3: errorPropellerSecurity: Propeller cutout security triggered
	//   4: errorCommLost: Communication with motor failed by timeout
	//   5: errorRCEmergencyStop: RC emergency stop
	//   6: errorRealTime: Motor controler scheduler real-time out of bounds
	//   7: errorMotorSetting: One or several incorrect values in motor settings
	//   8: errorTemperature: Too hot or too cold Cypress temperature
	//   9: errorBatteryVoltage: Battery voltage out of bounds
	//   10: errorLipoCells: Incorrect number of LIPO cells
	//   11: errorMOSFET: Defectuous MOSFET or broken motor phases
	//   12: errorBootloader: Not use for BLDC but useful for HAL
	//   13: errorAssert: Error Made by BLDC_ASSERT()
	log.Info("ardrone3.motorErrorStateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Motor version
// Description: Motor version.
// Support:
// Triggered:
// Result:
// WARNING: Deprecated
func (s *settingsState) motorSoftwareVersionChanged(args []interface{}) error {
	// version := args[0].(string)
	//   name of the version : dot separated fields (major version - minor version
	//   - firmware type - nb motors handled). Firmware types : Release, Debug,
	//   Alpha, Test-bench
	log.Info("ardrone3.motorSoftwareVersionChanged() called")
	return nil
}

// TODO: Implement this
// Title: Motor flight status
// Description: Motor flight status.
// Support: 0901;090c;090e
// Triggered: at connection.
// Result:
func (s *settingsState) motorFlightsStatusChanged(args []interface{}) error {
	// nbFlights := args[0].(uint16)
	//   total number of flights
	// lastFlightDuration := args[1].(uint16)
	//   Duration of the last flight (in seconds)
	// totalFlightDuration := args[2].(uint32)
	//   Duration of all flights (in seconds)
	log.Info("ardrone3.motorFlightsStatusChanged() called")
	return nil
}

// TODO: Implement this
// Title: Last motor error
// Description: Last motor error.\n This is a reminder of the last error. To
//   know if a motor error is currently happening, see [MotorError](#1-16-2).
// Support: 0901;090c;090e
// Triggered: at connection and when an error occurs.
// Result:
func (s *settingsState) motorErrorLastErrorChanged(args []interface{}) error {
	// motorError := args[0].(int32)
	//   Enumeration of the motor error
	//   0: noError: No error detected
	//   1: errorEEPRom: EEPROM access failure
	//   2: errorMotorStalled: Motor stalled
	//   3: errorPropellerSecurity: Propeller cutout security triggered
	//   4: errorCommLost: Communication with motor failed by timeout
	//   5: errorRCEmergencyStop: RC emergency stop
	//   6: errorRealTime: Motor controler scheduler real-time out of bounds
	//   7: errorMotorSetting: One or several incorrect values in motor settings
	//   8: errorBatteryVoltage: Battery voltage out of bounds
	//   9: errorLipoCells: Incorrect number of LIPO cells
	//   10: errorMOSFET: Defectuous MOSFET or broken motor phases
	//   11: errorTemperature: Too hot or too cold Cypress temperature
	//   12: errorBootloader: Not use for BLDC but useful for HAL
	//   13: errorAssert: Error Made by BLDC_ASSERT()
	log.Info("ardrone3.motorErrorLastErrorChanged() called")
	return nil
}

// p7ID is deprecated, but since we can still see this command being invoked,
// we'll implement the command to avoid a warning, but the implementation will
// remain a no-op unless / until such time that it becomes clear that older
// versions of the firmware might require us to support it.
func (s *settingsState) p7ID(args []interface{}) error {
	// serialID := args[0].(string)
	//   Product P7ID
	log.Debug("p7ID changed-- this is a no-op")
	return nil
}

// TODO: Implement this
func (s *settingsState) cPUID(args []interface{}) error {
	// id := args[0].(string)
	//   Product main cpu id
	log.Info("ardrone3.cPUID() called")
	return nil
}

func (s *settingsState) RLock() {
	s.lock.RLock()
}

func (s *settingsState) RUnlock() {
	s.lock.RUnlock()
}
