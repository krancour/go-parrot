package common

import (
	"sync"

	"github.com/krancour/go-parrot/lock"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Settings state from product

// SettingsState ...
// TODO: Document this
type SettingsState interface {
	lock.ReadLockable
	// ProductName returns the product's name. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	ProductName() (string, bool)
	// SerialHigh returns the high end of the device's serial number. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	SerialHigh() (string, bool)
	// SerialLow returns the low end of the device's serial number. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	SerialLow() (string, bool)
	// SoftwareVersion returns the device's software version. A boolean value is
	// also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	SoftwareVersion() (string, bool)
	// HardwareVersion returns the the device's hardware version. A boolean value
	// is also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	HardwareVersion() (string, bool)
	// AutoCountryEnabled returns a boolean indicating whether autoCountry is
	// enabled. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	AutoCountryEnabled() (bool, bool)
	// CountryCode returns the country code in ISO 3166 format. An empty string
	// indicates the country is unknown. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	CountryCode() (string, bool)
	// AllSettingsSent returns a boolean indicating whether the device has sent
	// all settings to the client. A boolean value is also returned, indicating
	// whether the first value was reported by the device (true) or a default
	// value (false). This permits callers to distinguish real zero values from
	// default zero values.
	AllSettingsSent() (bool, bool)
}

type settingsState struct {
	sync.RWMutex
	productName        *string
	serialHigh         *string
	serialLow          *string
	softwareVersion    *string
	hardwareVersion    *string
	autoCountryEnabled *bool
	countryCode        *string
	allSettingsSent    *bool
}

func newSettingsState() *settingsState {
	return &settingsState{}
}

func (s *settingsState) ClassID() uint8 {
	return 3
}

func (s *settingsState) ClassName() string {
	return "SettingsState"
}

func (s *settingsState) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"AllSettingsChanged",
			[]interface{}{},
			s.allSettingsChanged,
			log,
		),
		arcommands.NewD2CCommand(
			1,
			"ResetChanged",
			[]interface{}{},
			s.resetChanged,
			log,
		),
		arcommands.NewD2CCommand(
			2,
			"ProductNameChanged",
			[]interface{}{
				string(0), // name,
			},
			s.productNameChanged,
			log,
		),
		arcommands.NewD2CCommand(
			3,
			"ProductVersionChanged",
			[]interface{}{
				string(0), // software,
				string(0), // hardware,
			},
			s.productVersionChanged,
			log,
		),
		arcommands.NewD2CCommand(
			4,
			"ProductSerialHighChanged",
			[]interface{}{
				string(0), // high,
			},
			s.productSerialHighChanged,
			log,
		),
		arcommands.NewD2CCommand(
			5,
			"ProductSerialLowChanged",
			[]interface{}{
				string(0), // low,
			},
			s.productSerialLowChanged,
			log,
		),
		arcommands.NewD2CCommand(
			6,
			"CountryChanged",
			[]interface{}{
				string(0), // code,
			},
			s.countryChanged,
			log,
		),
		arcommands.NewD2CCommand(
			7,
			"AutoCountryChanged",
			[]interface{}{
				uint8(0), // automatic,
			},
			s.autoCountryChanged,
			log,
		),
	}
}

// Invoked by the device to indicate all settings have been sent.
func (s *settingsState) allSettingsChanged(
	args []interface{},
	log *log.Entry,
) error {
	s.Lock()
	defer s.Unlock()
	s.allSettingsSent = ptr.ToBool(true)
	log.Debug("all settings have been sent by the device")
	return nil
}

// TODO: Implement this
// Title: All settings have been reset
// Description: All settings have been reset.
// Support: drones
// Triggered: by [ResetSettings](#0-2-1).
// Result:
func (s *settingsState) resetChanged(
	args []interface{},
	log *log.Entry,
) error {
	log.Warn("command not implemented")
	return nil
}

// productNameChanged is invoked by the device to indicate its name has changed.
func (s *settingsState) productNameChanged(
	args []interface{},
	log *log.Entry,
) error {
	s.Lock()
	defer s.Unlock()
	s.productName = ptr.ToString(args[0].(string))
	log.WithField(
		"productName", *s.productName,
	).Debug("product name changed")
	return nil
}

// productVersionChanged is invoked during the connection process to report
// the product version.
func (s *settingsState) productVersionChanged(
	args []interface{},
	log *log.Entry,
) error {
	s.Lock()
	defer s.Unlock()
	s.softwareVersion = ptr.ToString(args[0].(string))
	s.hardwareVersion = ptr.ToString(args[1].(string))
	log.WithField(
		"softwareVersion", *s.softwareVersion,
	).WithField(
		"hardwareVersion", *s.hardwareVersion,
	).Debug("common.productVersionChanged() called")
	return nil
}

// productSerialHighChanged is invoked during the connection process to report
// the the high end of the device's serial number.
func (s *settingsState) productSerialHighChanged(
	args []interface{},
	log *log.Entry,
) error {
	s.Lock()
	defer s.Unlock()
	s.serialHigh = ptr.ToString(args[0].(string))
	log.WithField(
		"serialHigh", *s.serialHigh,
	).Debug("serial high changed")
	return nil
}

// productSerialLowChanged is invoked during the connection process to report
// the the low end of the device's serial number.
func (s *settingsState) productSerialLowChanged(
	args []interface{},
	log *log.Entry,
) error {
	s.Lock()
	defer s.Unlock()
	s.serialLow = ptr.ToString(args[0].(string))
	log.WithField(
		"serialLow", *s.serialLow,
	).Debug("serial low changed")
	return nil
}

// countryChanged is invoked by the device when its country is changed.
func (s *settingsState) countryChanged(
	args []interface{},
	log *log.Entry,
) error {
	s.Lock()
	defer s.Unlock()
	s.countryCode = ptr.ToString(args[0].(string))
	//   Country code with ISO 3166 format, empty string means unknown country.
	log.WithField(
		"countryCode", *s.countryCode,
	).Debug("common.countryChanged() called")
	return nil
}

// autoCountryChanged is invoked by the device to indicate whether auto country
// is enabled.
func (s *settingsState) autoCountryChanged(
	args []interface{},
	log *log.Entry,
) error {
	s.Lock()
	defer s.Unlock()
	s.autoCountryEnabled = ptr.ToBool(args[0].(uint8) == 1)
	log.WithField(
		"autoCountryEnabled", *s.autoCountryEnabled,
	).Debug("autoCountry changed")
	return nil
}

func (s *settingsState) ProductName() (string, bool) {
	if s.productName == nil {
		return "", false
	}
	return *s.productName, true
}

func (s *settingsState) SerialHigh() (string, bool) {
	if s.serialHigh == nil {
		return "", false
	}
	return *s.serialHigh, true
}

func (s *settingsState) SerialLow() (string, bool) {
	if s.serialLow == nil {
		return "", false
	}
	return *s.serialLow, true
}

func (s *settingsState) SoftwareVersion() (string, bool) {
	if s.softwareVersion == nil {
		return "", false
	}
	return *s.softwareVersion, true
}

func (s *settingsState) HardwareVersion() (string, bool) {
	if s.hardwareVersion == nil {
		return "", false
	}
	return *s.hardwareVersion, true
}

func (s *settingsState) AutoCountryEnabled() (bool, bool) {
	if s.autoCountryEnabled == nil {
		return false, false
	}
	return *s.autoCountryEnabled, true
}

func (s *settingsState) CountryCode() (string, bool) {
	if s.countryCode == nil {
		return "", false
	}
	return *s.countryCode, true
}

func (s *settingsState) AllSettingsSent() (bool, bool) {
	if s.allSettingsSent == nil {
		return false, false
	}
	return *s.allSettingsSent, true
}
