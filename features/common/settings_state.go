package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Settings state from product

// SettingsState ...
// TODO: Document this
type SettingsState interface{}

type settingsState struct{}

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

// TODO: Implement this
// Title: Product name changed
// Description: Product name changed.
// Support: drones
// Triggered: by [SetProductName](#0-2-2).
// Result:
func (s *settingsState) productNameChanged(args []interface{}) error {
	// name := args[0].(string)
	//   Product name
	log.Info("common.productNameChanged() called")
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
