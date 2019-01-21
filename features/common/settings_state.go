package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Settings state from product

// SettingsState ...
// TODO: Document this
type SettingsState interface {
	// RLock blocks until a read lock is obtained. This permits callers to procede
	// with querying any or all attributes of the settings state without worry
	// that some attributes will be overwritten as others are read. i.e. It
	// permits the possibility of taking an atomic snapshop of settings state.
	// Note that use of this function is not obligatory for applications that do
	// not require such guarantees. Callers MUST call RUnlock() or else settings
	// state will never resume updating.
	RLock()
	// RUnlock releases a read lock on the GPS state. See RLock().
	RUnlock()
	// ProductName returns the product's name. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	ProductName() (string, bool)
}

type settingsState struct {
	productName *string
	lock        sync.RWMutex
}

func (s *settingsState) ID() uint8 {
	return 3
}

func (s *settingsState) Name() string {
	return "SettingsState"
}

func (s *settingsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"AllSettingsChanged",
			[]interface{}{},
			s.allSettingsChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"ResetChanged",
			[]interface{}{},
			s.resetChanged,
		),
		arcommands.NewD2CCommand(
			2,
			"ProductNameChanged",
			[]interface{}{
				string(0), // name,
			},
			s.productNameChanged,
		),
		arcommands.NewD2CCommand(
			3,
			"ProductVersionChanged",
			[]interface{}{
				string(0), // software,
				string(0), // hardware,
			},
			s.productVersionChanged,
		),
		arcommands.NewD2CCommand(
			4,
			"ProductSerialHighChanged",
			[]interface{}{
				string(0), // high,
			},
			s.productSerialHighChanged,
		),
		arcommands.NewD2CCommand(
			5,
			"ProductSerialLowChanged",
			[]interface{}{
				string(0), // low,
			},
			s.productSerialLowChanged,
		),
		arcommands.NewD2CCommand(
			6,
			"CountryChanged",
			[]interface{}{
				string(0), // code,
			},
			s.countryChanged,
		),
		arcommands.NewD2CCommand(
			7,
			"AutoCountryChanged",
			[]interface{}{
				uint8(0), // automatic,
			},
			s.autoCountryChanged,
		),
	}
}

// TODO: Implement this
// Title: All settings have been sent
// Description: All settings have been sent.\n\n **Please note that you should
//   not care about this event if you are using the libARController API as this
//   library is handling the connection process for you.**
// Support: drones
// Triggered: when all settings values have been sent.
// Result:
func (s *settingsState) allSettingsChanged(args []interface{}) error {
	log.Info("common.allSettingsChanged() called")
	return nil
}

// TODO: Implement this
// Title: All settings have been reset
// Description: All settings have been reset.
// Support: drones
// Triggered: by [ResetSettings](#0-2-1).
// Result:
func (s *settingsState) resetChanged(args []interface{}) error {
	log.Info("common.resetChanged() called")
	return nil
}

// productNameChanged is invoked by the device to indicate its name has changed.
func (s *settingsState) productNameChanged(args []interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.productName = ptr.ToString(args[0].(string))
	log.WithField(
		"productName", *s.productName,
	).Debug("product name changed")
	return nil
}

// TODO: Implement this
// Title: Product version
// Description: Product version.
// Support: drones
// Triggered: during the connection process.
// Result:
func (s *settingsState) productVersionChanged(args []interface{}) error {
	// software := args[0].(string)
	//   Product software version
	// hardware := args[1].(string)
	//   Product hardware version
	log.Info("common.productVersionChanged() called")
	return nil
}

// TODO: Implement this
// Title: Product serial (1st part)
// Description: Product serial (1st part).
// Support: drones
// Triggered: during the connection process.
// Result:
func (s *settingsState) productSerialHighChanged(args []interface{}) error {
	// high := args[0].(string)
	//   Serial high number (hexadecimal value)
	log.Info("common.productSerialHighChanged() called")
	return nil
}

// TODO: Implement this
// Title: Product serial (2nd part)
// Description: Product serial (2nd part).
// Support: drones
// Triggered: during the connection process.
// Result:
func (s *settingsState) productSerialLowChanged(args []interface{}) error {
	// low := args[0].(string)
	//   Serial low number (hexadecimal value)
	log.Info("common.productSerialLowChanged() called")
	return nil
}

// TODO: Implement this
// Title: Country changed
// Description: Country changed.
// Support: drones
// Triggered: by [SetCountry](#0-2-3).
// Result:
func (s *settingsState) countryChanged(args []interface{}) error {
	// code := args[0].(string)
	//   Country code with ISO 3166 format, empty string means unknown country.
	log.Info("common.countryChanged() called")
	return nil
}

// TODO: Implement this
// Title: Auto-country changed
// Description: Auto-country changed.
// Support: drones
// Triggered: by [SetAutoCountry](#0-2-4).
// Result:
func (s *settingsState) autoCountryChanged(args []interface{}) error {
	// automatic := args[0].(uint8)
	//   Boolean : 0 : Manual / 1 : Auto
	log.Info("common.autoCountryChanged() called")
	return nil
}

func (s *settingsState) RLock() {
	s.lock.RLock()
}

func (s *settingsState) RUnlock() {
	s.lock.RUnlock()
}

func (s *settingsState) ProductName() (string, bool) {
	if s.productName == nil {
		return "", false
	}
	return *s.productName, true
}
