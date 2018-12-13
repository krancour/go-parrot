package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Commands sent by the firmware to advertise the charger status.

// ChargerState ...
// TODO: Document this
type ChargerState interface{}

type chargerState struct{}

func (c *chargerState) ID() uint8 {
	return 29
}

func (c *chargerState) Name() string {
	return "ChargerState"
}

func (c *chargerState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"MaxChargeRateChanged",
			[]interface{}{
				int32(0), // rate,
			},
			c.maxChargeRateChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"CurrentChargeStateChanged",
			[]interface{}{
				int32(0), // status,
				int32(0), // phase,
			},
			c.currentChargeStateChanged,
		),
		arcommands.NewD2CCommand(
			2,
			"LastChargeRateChanged",
			[]interface{}{
				int32(0), // rate,
			},
			c.lastChargeRateChanged,
		),
		arcommands.NewD2CCommand(
			3,
			"ChargingInfo",
			[]interface{}{
				int32(0), // phase,
				int32(0), // rate,
				uint8(0), // intensity,
				uint8(0), // fullChargingTime,
			},
			c.chargingInfo,
		),
	}
}

// TODO: Implement this
// Title: Max charge rate
// Description: Max charge rate.
// Support:
// Triggered:
// Result:
// WARNING: Deprecated
func (c *chargerState) maxChargeRateChanged(args []interface{}) error {
	// rate := args[0].(int32)
	//   The current maximum charge rate.
	//   0: SLOW: Fully charge the battery at a slow rate. Typically limit max
	//      charge current to 512 mA.
	//   1: MODERATE: Almost fully-charge the battery at moderate rate
	//      (&gt; 512 mA) but slower than the fastest rate.
	//   2: FAST: Almost fully-charge the battery at the highest possible rate
	//      supported by the charger.
	log.Info("common.maxChargeRateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Current charge state
// Description: Current charge state.
// Support:
// Triggered:
// Result:
// WARNING: Deprecated
func (c *chargerState) currentChargeStateChanged(args []interface{}) error {
	// status := args[0].(int32)
	//   Charger status.
	//   0: DISCHARGING: The battery is discharging.
	//   1: CHARGING_SLOW: The battery is charging at a slow rate about 512 mA.
	//   2: CHARGING_MODERATE: The battery is charging at a moderate rate
	//      (&gt; 512 mA) but slower than the fastest rate.
	//   3: CHARGING_FAST: The battery is charging at a the fastest rate.
	//   4: BATTERY_FULL: The charger is plugged and the battery is fully charged.
	// phase := args[1].(int32)
	//   The current charging phase.
	//   0: UNKNOWN: The charge phase is unknown or irrelevant.
	//   1: CONSTANT_CURRENT_1: First phase of the charging process. The battery
	//      is charging with constant current.
	//   2: CONSTANT_CURRENT_2: Second phase of the charging process. The battery
	//      is charging with constant current, with a higher voltage than the
	//      first phase.
	//   3: CONSTANT_VOLTAGE: Last part of the charging process. The battery is
	//      charging with a constant voltage.
	//   4: CHARGED: The battery is fully charged.
	log.Info("common.currentChargeStateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Last charge rate
// Description: Last charge rate.
// Support:
// Triggered:
// Result:
// WARNING: Deprecated
func (c *chargerState) lastChargeRateChanged(args []interface{}) error {
	// rate := args[0].(int32)
	//   The charge rate recorded by the firmware for the last charge.
	//   0: UNKNOWN: The last charge rate is not known.
	//   1: SLOW: Slow charge rate.
	//   2: MODERATE: Moderate charge rate.
	//   3: FAST: Fast charge rate.
	log.Info("common.lastChargeRateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Charging information
// Description: Charging information.
// Support: 0905;0906;0907;0909;090a
// Triggered: when the product is charging or when the charging state changes.
// Result:
func (c *chargerState) chargingInfo(args []interface{}) error {
	// phase := args[0].(int32)
	//   The current charging phase.
	//   0: UNKNOWN: The charge phase is unknown or irrelevant.
	//   1: CONSTANT_CURRENT_1: First phase of the charging process. The battery
	//      is charging with constant current.
	//   2: CONSTANT_CURRENT_2: Second phase of the charging process. The battery
	//      is charging with constant current, with a higher voltage than the
	//      first phase.
	//   3: CONSTANT_VOLTAGE: Last part of the charging process. The battery is
	//      charging with a constant voltage.
	//   4: CHARGED: The battery is fully charged.
	//   5: DISCHARGING: The battery is discharging; Other arguments refers to the
	//      last charge.
	// rate := args[1].(int32)
	//   The charge rate. If phase is DISCHARGING, refers to the last charge.
	//   0: UNKNOWN: The charge rate is not known.
	//   1: SLOW: Slow charge rate.
	//   2: MODERATE: Moderate charge rate.
	//   3: FAST: Fast charge rate.
	// intensity := args[2].(uint8)
	//   The charging intensity, in dA. (12dA = 1,2A) ; If phase is DISCHARGING,
	//   refers to the last charge. Equals to 0 if not known.
	// fullChargingTime := args[3].(uint8)
	//   The full charging time estimated, in minute. If phase is DISCHARGING,
	//   refers to the last charge. Equals to 0 if not known.
	log.Info("common.chargingInfo() called")
	return nil
}
