package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Status of the calibration

// CalibrationState ...
// TODO: Document this
type CalibrationState interface{}

type calibrationState struct{}

func (c *calibrationState) ID() uint8 {
	return 14
}

func (c *calibrationState) Name() string {
	return "CalibrationState"
}

func (c *calibrationState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"MagnetoCalibrationStateChanged",
			[]interface{}{
				uint8(0), // xAxisCalibration,
				uint8(0), // yAxisCalibration,
				uint8(0), // zAxisCalibration,
				uint8(0), // calibrationFailed,
			},
			c.magnetoCalibrationStateChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"MagnetoCalibrationRequiredState",
			[]interface{}{
				uint8(0), // required,
			},
			c.magnetoCalibrationRequiredState,
		),
		arcommands.NewD2CCommand(
			2,
			"MagnetoCalibrationAxisToCalibrateChanged",
			[]interface{}{
				int32(0), // axis,
			},
			c.magnetoCalibrationAxisToCalibrateChanged,
		),
		arcommands.NewD2CCommand(
			3,
			"MagnetoCalibrationStartedChanged",
			[]interface{}{
				uint8(0), // started,
			},
			c.magnetoCalibrationStartedChanged,
		),
		arcommands.NewD2CCommand(
			4,
			"PitotCalibrationStateChanged",
			[]interface{}{
				int32(0), // state,
				uint8(0), // lastError,
			},
			c.pitotCalibrationStateChanged,
		),
	}
}

// TODO: Implement this
// Title: Magneto calib process axis state
// Description: Magneto calib process axis state.
// Support: 0901;090c;090e
// Triggered: when the calibration process is started with
//   [StartOrAbortMagnetoCalib](#0-13-0) and each time an axis calibration state
//   changes.
// Result:
func (c *calibrationState) magnetoCalibrationStateChanged(
	args []interface{},
) error {
	// xAxisCalibration := args[0].(uint8)
	//   State of the x axis (roll) calibration : 1 if calibration is done, 0
	//   otherwise
	// yAxisCalibration := args[1].(uint8)
	//   State of the y axis (pitch) calibration : 1 if calibration is done, 0
	//   otherwise
	// zAxisCalibration := args[2].(uint8)
	//   State of the z axis (yaw) calibration : 1 if calibration is done, 0
	//   otherwise
	// calibrationFailed := args[3].(uint8)
	//   1 if calibration has failed, 0 otherwise. If this arg is 1, consider all
	//   previous arg as 0
	log.Info("common.magnetoCalibrationStateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Calibration required
// Description: Calibration required.
// Support: 0901;090c;090e
// Triggered: when the calibration requirement changes.
// Result:
func (c *calibrationState) magnetoCalibrationRequiredState(
	args []interface{},
) error {
	// required := args[0].(uint8)
	//   1 if calibration is required, 0 if current calibration is still valid
	log.Info("common.magnetoCalibrationRequiredState() called")
	return nil
}

// TODO: Implement this
// Title: Axis to calibrate during calibration process
// Description: Axis to calibrate during calibration process.
// Support: 0901;090c;090e
// Triggered: during the calibration process when the axis to calibrate changes.
// Result:
func (c *calibrationState) magnetoCalibrationAxisToCalibrateChanged(
	args []interface{},
) error {
	// axis := args[0].(int32)
	//   The axis to calibrate
	//   0: xAxis: If the current calibration axis should be the x axis
	//   1: yAxis: If the current calibration axis should be the y axis
	//   2: zAxis: If the current calibration axis should be the z axis
	//   3: none: If none of the axis should be calibrated
	log.Info("common.magnetoCalibrationAxisToCalibrateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Calibration process state
// Description: Calibration process state.
// Support: 0901;090c;090e
// Triggered: by [StartOrAbortMagnetoCalib](#0-13-0) or when the process ends
//   because it succeeded.
// Result:
func (c *calibrationState) magnetoCalibrationStartedChanged(
	args []interface{},
) error {
	// started := args[0].(uint8)
	//   1 if calibration has started, 0 otherwise
	log.Info("common.magnetoCalibrationStartedChanged() called")
	return nil
}

// TODO: Implement this
func (c *calibrationState) pitotCalibrationStateChanged(
	args []interface{},
) error {
	// state := args[0].(int32)
	//   State of pitot calibration
	//   0: done: Calibration is ok
	//   1: ready: Calibration is started, waiting user action
	//   2: in_progress: Calibration is in progress
	//   3: required: Calibration is required
	// lastError := args[1].(uint8)
	//   lastError : 1 if an error occured and 0 if not
	log.Info("common.pitotCalibrationStateChanged() called")
	return nil
}
