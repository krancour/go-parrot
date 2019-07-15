package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Settings state from product

const (
	MotorFrontLeft  uint8 = 1
	MotorFrontRight uint8 = 2
	MotorBackRight  uint8 = 4
	MotorBackLeft   uint8 = 8
)

// SettingsState ...
// TODO: Document this
type SettingsState interface {
	lock.ReadLockable
	// GPSSoftwareVersion returns the software version of the product's GPS. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	GPSSoftwareVersion() (string, bool)
	// GPSHardwareVersion hardware version of the product's GPS. A boolean value
	// is also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	GPSHardwareVersion() (string, bool)
	// NumberOfFlights returns the total number of flights completed. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	NumberOfFlights() (uint16, bool)
	// LastFlightDuration returns the length, in seconds, of the last flight. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	LastFlightDuration() (uint16, bool)
	// TotalFlightDuration returns the cummulative length of all completed flights
	// in seconds. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	TotalFlightDuration() (uint32, bool)
	// CPUID returns the drone's CPU's ID. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	CPUID() (string, bool)
}

type settingsState struct {
	sync.RWMutex
	// gpsSoftwareVersion represents the software version of the product's GPS
	gpsSoftwareVersion *string
	// gpsHardwareVersion represents the hardware version of the product's GPS
	gpsHardwareVersion *string
	// numberOfFlights represents the total number of flights completed by this
	// drone
	numberOfFlights *uint16
	// lastFlightDuration represents the length, in seconds, of the last flight
	lastFlightDuration *uint16
	// totalFlightDuration represents the cummulative length of all completed
	// flights in seconds
	totalFlightDuration *uint32
	// deviceCPUID represents the drone's CPU's ID
	deviceCPUID *string
}

func (s *settingsState) ID() uint8 {
	return 16
}

func (s *settingsState) Name() string {
	return "SettingsState"
}

func (s *settingsState) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			1,
			"ProductGPSVersionChanged",
			[]interface{}{
				string(0), // software,
				string(0), // hardware,
			},
			s.productGPSVersionChanged,
			log,
		),
		arcommands.NewD2CCommand(
			2,
			"MotorErrorStateChanged",
			[]interface{}{
				uint8(0), // motorIds,
				int32(0), // motorError,
			},
			s.motorErrorStateChanged,
			log,
		),
		arcommands.NewD2CCommand(
			3,
			"MotorSoftwareVersionChanged",
			[]interface{}{
				string(0), // version,
			},
			s.motorSoftwareVersionChanged,
			log,
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
			log,
		),
		arcommands.NewD2CCommand(
			5,
			"MotorErrorLastErrorChanged",
			[]interface{}{
				int32(0), // motorError,
			},
			s.motorErrorLastErrorChanged,
			log,
		),
		arcommands.NewD2CCommand(
			6,
			"P7ID",
			[]interface{}{
				string(0), // serialID,
			},
			s.p7ID,
			log,
		),
		arcommands.NewD2CCommand(
			7,
			"CPUID",
			[]interface{}{
				string(0), // id,
			},
			s.cpuID,
			log,
		),
	}
}

// productGPSVersionChanged is invoked by the device at connection time.
func (s *settingsState) productGPSVersionChanged(
	args []interface{},
	log *log.Entry,
) error {
	s.Lock()
	defer s.Unlock()
	s.gpsSoftwareVersion = ptr.ToString(args[0].(string))
	s.gpsHardwareVersion = ptr.ToString(args[1].(string))
	log.WithField(
		"software", *s.gpsSoftwareVersion,
	).WithField(
		"hardware", *s.gpsHardwareVersion,
	).Debug("product gps version changed")
	return nil
}

// motorErrorStateChanged is triggered by the device when a motor error occurs.
func (s *settingsState) motorErrorStateChanged(
	args []interface{},
	log *log.Entry,
) error {
	motorIds := args[0].(uint8)
	l := log.WithField(
		"fontLeft", motorIds&MotorFrontLeft == MotorFrontLeft,
	).WithField(
		"fontRight", motorIds&MotorFrontRight == MotorFrontRight,
	).WithField(
		"backRight", motorIds&MotorBackRight == MotorBackRight,
	).WithField(
		"backLeft", motorIds&MotorBackLeft == MotorBackLeft,
	)
	motorError := args[1].(int32)
	switch motorError {
	case 0:
		log.Debug("no motor errors detected")
	case 1:
		l.Error("motor error: EEPROM access failure")
	case 2:
		l.Error("motor error: errorMotorStalled: Motor stalled")
	case 3:
		l.Error("motor error: errorPropellerSecurity: Propeller cutout security triggered")
	case 4:
		l.Error("motor error: errorCommLost: Communication with motor failed by timeout")
	case 5:
		l.Error("motor error: errorRCEmergencyStop: RC emergency stop")
	case 6:
		l.Error("motor error: errorRealTime: Motor controler scheduler real-time out of bounds")
	case 7:
		l.Error("motor error: errorMotorSetting: One or several incorrect values in motor settings")
	case 8:
		l.Error("motor error: errorTemperature: Too hot or too cold Cypress temperature")
	case 9:
		l.Error("motor error: errorBatteryVoltage: Battery voltage out of bounds")
	case 10:
		l.Error("motor error: errorLipoCells: Incorrect number of LIPO cells")
	case 11:
		l.Error("motor error: errorMOSFET: Defectuous MOSFET or broken motor phases")
	case 12:
		l.Error("motor error: errorBootloader: Not use for BLDC but useful for HAL")
	case 13:
		l.Error("motor error: errorAssert: Error Made by BLDC_ASSERT()")
	}
	return nil
}

// motorSoftwareVersionChanged is deprecated, but since we can still see this
// command being invoked, we'll implement the command to avoid a warning, but
// the implementation will remain a no-op unless / until such time that it
// becomes clear that older versions of the firmware might require us to support
// it.
func (s *settingsState) motorSoftwareVersionChanged(
	args []interface{},
	log *log.Entry,
) error {
	// version := args[0].(string)
	//   name of the version : dot separated fields (major version - minor version
	//   - firmware type - nb motors handled). Firmware types : Release, Debug,
	//   Alpha, Test-bench
	log.Debug("motor software version changed-- this is a no-op")
	return nil
}

// motorFlightsStatusChanged is invoked by the device at connection time.
func (s *settingsState) motorFlightsStatusChanged(
	args []interface{},
	log *log.Entry,
) error {
	s.Lock()
	defer s.Unlock()
	s.numberOfFlights = ptr.ToUint16(args[0].(uint16))
	s.lastFlightDuration = ptr.ToUint16(args[1].(uint16))
	s.totalFlightDuration = ptr.ToUint32(args[2].(uint32))
	log.WithField(
		"nbFlights", *s.numberOfFlights,
	).WithField(
		"lastFlightDuration", *s.lastFlightDuration,
	).WithField(
		"totalFlightDuration", *s.totalFlightDuration,
	).Debug("motor flight status changed")
	return nil
}

// motorErrorLastErrorChanged is a no-op because the function it's meant to
// carry out-- a "reminder" of the most recent motor error is not contextually
// necessary for this SDK.
func (s *settingsState) motorErrorLastErrorChanged(
	args []interface{},
	log *log.Entry,
) error {
	log.Debug("last motor error changed-- this is a no-op")
	return nil
}

// p7ID is deprecated, but since we can still see this command being invoked,
// we'll implement the command to avoid a warning, but the implementation will
// remain a no-op unless / until such time that it becomes clear that older
// versions of the firmware might require us to support it.
func (s *settingsState) p7ID(args []interface{}, log *log.Entry) error {
	// serialID := args[0].(string)
	//   Product P7ID
	log.Debug("p7ID changed-- this is a no-op")
	return nil
}

// cpuID is invoked by the device at connecton time to report its CPU ID.
func (s *settingsState) cpuID(args []interface{}, log *log.Entry) error {
	s.Lock()
	defer s.Unlock()
	s.deviceCPUID = ptr.ToString(args[0].(string))
	log.WithField(
		"id", *s.deviceCPUID,
	).Debug("device CPU ID changed")
	return nil
}

func (s *settingsState) GPSSoftwareVersion() (string, bool) {
	if s.gpsSoftwareVersion == nil {
		return "", false
	}
	return *s.gpsSoftwareVersion, true
}

func (s *settingsState) GPSHardwareVersion() (string, bool) {
	if s.gpsHardwareVersion == nil {
		return "", false
	}
	return *s.gpsHardwareVersion, true
}

func (s *settingsState) NumberOfFlights() (uint16, bool) {
	if s.numberOfFlights == nil {
		return 0, false
	}
	return *s.numberOfFlights, true
}

func (s *settingsState) LastFlightDuration() (uint16, bool) {
	if s.lastFlightDuration == nil {
		return 0, false
	}
	return *s.lastFlightDuration, true
}

func (s *settingsState) TotalFlightDuration() (uint32, bool) {
	if s.totalFlightDuration == nil {
		return 0, false
	}
	return *s.totalFlightDuration, true
}

func (s *settingsState) CPUID() (string, bool) {
	if s.deviceCPUID == nil {
		return "", false
	}
	return *s.deviceCPUID, true
}
