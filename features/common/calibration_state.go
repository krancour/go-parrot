package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Status of the calibration

const (
	MagnetoCalibrationAxisX    int32 = 0
	MagnetoCalibrationAxisY    int32 = 1
	MagnetoCalibrationAxisZ    int32 = 2
	MagnetoCalibrationAxisNone int32 = 3
)

// CalibrationState ...
// TODO: Document this
type CalibrationState interface {
	// RLock blocks until a read lock is obtained. This permits callers to procede
	// with querying any or all attributes of the calibration state without worry
	// that some attributes will be overwritten as others are read. i.e. It
	// permits the possibility of taking an atomic snapshop of calibration state.
	// Note that use of this function is not obligatory for applications that do
	// not require such guarantees. Callers MUST call RUnlock() or else
	// calibration state will never resume updating.
	RLock()
	// RUnlock releases a read lock on the calibration state. See RLock().
	RUnlock()
	// MagnetoCalibrationRequired returns a boolean indicating whether the device
	// requires magneto calibrartion to be performed. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	MagnetoCalibrationRequired() (bool, bool)
	// MagnetoCalibrationStarted returns a boolean indicating whether magneto
	// calibration is currently in progress. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	MagnetoCalibrationStarted() (bool, bool)
	// MagnetoCalibrationAxis returns an int32 indicating which axis, if any is
	// currently being calibrated. A boolean value is also returned, indicating
	// whether the first value was reported by the device (true) or a default
	// value (false). This permits callers to distinguish real zero values from
	// default zero values.
	MagnetoCalibrationAxis() (int32, bool)
}

type calibrationState struct {
	magnetoCalibrationRequired *bool
	magnetoCalibrationStarted  *bool
	magnetoCalibrationAxis     *int32
	lock                       sync.RWMutex
}

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

// magnetoCalibrationRequiredState is invoked by the device to indicate that
// magnetometer calibration is required.
func (c *calibrationState) magnetoCalibrationRequiredState(
	args []interface{},
) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.magnetoCalibrationRequired = ptr.ToBool(args[0].(uint8) == 1)
	log.WithField(
		"magnetoCalibrationRequired", *c.magnetoCalibrationRequired,
	).Debug("magneto claibration required state changed")
	return nil
}

// magnetoCalibrationAxisToCalibrateChanged is invoked by the device to indicate
// which axis is actively being calibrated.
func (c *calibrationState) magnetoCalibrationAxisToCalibrateChanged(
	args []interface{},
) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.magnetoCalibrationAxis = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"magnetoCalibrationAxis", *c.magnetoCalibrationAxis,
	).Debug("common.magnetoCalibrationAxisToCalibrateChanged() called")
	return nil
}

// magnetoCalibrationStartedChanged is invoked by the device to indicate whether
// magneto calibration is in progress.
func (c *calibrationState) magnetoCalibrationStartedChanged(
	args []interface{},
) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.magnetoCalibrationStarted = ptr.ToBool(args[0].(uint8) == 1)
	log.WithField(
		"magnetoCalibrationStarted", *c.magnetoCalibrationStarted,
	).Debug("magneto calibration started state changed")
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

func (c *calibrationState) RLock() {
	c.lock.RLock()
}

func (c *calibrationState) RUnlock() {
	c.lock.RUnlock()
}

func (c *calibrationState) MagnetoCalibrationRequired() (bool, bool) {
	if c.magnetoCalibrationRequired == nil {
		return false, false
	}
	return *c.magnetoCalibrationRequired, true
}

func (c *calibrationState) MagnetoCalibrationStarted() (bool, bool) {
	if c.magnetoCalibrationStarted == nil {
		return false, false
	}
	return *c.magnetoCalibrationStarted, true
}

func (c *calibrationState) MagnetoCalibrationAxis() (int32, bool) {
	if c.magnetoCalibrationAxis == nil {
		return 0, false
	}
	return *c.magnetoCalibrationAxis, true
}
