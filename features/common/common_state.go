package common

import (
	"sync"

	"github.com/krancour/go-parrot/lock"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// TODO: Mass storage stuff all needs to be in a map indexed by ID

// Common state from product

// CommonState ...
// TODO: Document this
type CommonState interface {
	lock.ReadLockable
	// RSSI returns the relative signal stength between the client and the device
	// in dbm. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	RSSI() (int16, bool)
	// BatteryPercent returns the percentage of battery life remaining. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	BatteryPercent() (uint8, bool)
	// IMUOK returns a boolean indicating whether the inertial measurement unit
	// sensor is operating normally. A boolean value is also returned, indicating
	// whether the first value was reported by the device (true) or a default
	// value (false). This permits callers to distinguish real zero values from
	// default zero values.
	IMUOK() (bool, bool)
	// BarometerOK returns a boolean indicating whether the barometer is
	// operarting normally. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	BarometerOK() (bool, bool)
	// UltrasoundOK returns a boolean indicating whether the ultrasonic sensor is
	// operating normally. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	UltrasoundOK() (bool, bool)
	// GPSOK returns a boolean indicating whether the GPS is operating normally. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	GPSOK() (bool, bool)
	// MagnetometerOK returns a boolean indicating whether the magnetometer is
	// operating normally. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	MagnetometerOK() (bool, bool)
	// VerticalCameraOK returns a boolean indicating whether the vertical camera
	// is operating normally. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	VerticalCameraOK() (bool, bool)
	// AllStatesSent returns a boolean indicating whether the device has sent
	// all states to the client. A boolean value is also returned, indicating
	// whether the first value was reported by the device (true) or a default
	// value (false). This permits callers to distinguish real zero values from
	// default zero values.
	AllStatesSent() (bool, bool)
	// MassStorageDevices returns a map of MassStorageDevices indexed by ID.
	MassStorageDevices() map[uint8]MassStorageDevice
}

type commonState struct {
	sync.RWMutex
	// TODO: Is this right? I thought RSSI is a relative measure, while dbm
	// would seem to indicate an absolute measure.
	// rssi is the relative signal stength between the client and the device
	// in dbm
	rssi               *int16
	massStorageDevices map[uint8]MassStorageDevice
	batteryPercent     *uint8
	imuOK              *bool
	barometerOK        *bool
	ultrasoundOK       *bool
	gpsOK              *bool
	magnetomenterOK    *bool
	verticalCameraOK   *bool
	allStatesSent      *bool
}

func (c *commonState) ID() uint8 {
	return 5
}

func (c *commonState) Name() string {
	return "CommonState"
}

func (c *commonState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"AllStatesChanged",
			[]interface{}{},
			c.allStatesChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"BatteryStateChanged",
			[]interface{}{
				uint8(0), // percent,
			},
			c.batteryStateChanged,
		),
		arcommands.NewD2CCommand(
			2,
			"MassStorageStateListChanged",
			[]interface{}{
				uint8(0),  // mass_storage_id,
				string(0), // name,
			},
			c.massStorageStateListChanged,
		),
		arcommands.NewD2CCommand(
			3,
			"MassStorageInfoStateListChanged",
			[]interface{}{
				uint8(0),  // mass_storage_id,
				uint32(0), // size,
				uint32(0), // used_size,
				uint8(0),  // plugged,
				uint8(0),  // full,
				uint8(0),  // internal,
			},
			c.massStorageInfoStateListChanged,
		),
		arcommands.NewD2CCommand(
			4,
			"CurrentDateChanged",
			[]interface{}{
				string(0), // date,
			},
			c.currentDateChanged,
		),
		arcommands.NewD2CCommand(
			5,
			"CurrentTimeChanged",
			[]interface{}{
				string(0), // time,
			},
			c.currentTimeChanged,
		),
		// arcommands.NewD2CCommand(
		// 	6,
		// 	"MassStorageInfoRemainingListChanged",
		// 	[]interface{}{
		// 		uint32(0), // free_space,
		// 		uint16(0), // rec_time,
		// 		uint32(0), // photo_remaining,
		// 	},
		// 	c.massStorageInfoRemainingListChanged,
		// ),
		arcommands.NewD2CCommand(
			7,
			"WifiSignalChanged",
			[]interface{}{
				int16(0), // rssi,
			},
			c.wifiSignalChanged,
		),
		arcommands.NewD2CCommand(
			8,
			"SensorsStatesListChanged",
			[]interface{}{
				int32(0), // sensorName,
				uint8(0), // sensorState,
			},
			c.sensorsStatesListChanged,
		),
		// arcommands.NewD2CCommand(
		// 	9,
		// 	"ProductModel",
		// 	[]interface{}{
		// 		int32(0), // model,
		// 	},
		// 	c.productModel,
		// ),
		// arcommands.NewD2CCommand(
		// 	10,
		// 	"CountryListKnown",
		// 	[]interface{}{
		// 		uint8(0),  // listFlags,
		// 		string(0), // countryCodes,
		// 	},
		// 	c.countryListKnown,
		// ),
		arcommands.NewD2CCommand(
			11,
			"DeprecatedMassStorageContentChanged",
			[]interface{}{
				uint8(0),  // mass_storage_id,
				uint16(0), // nbPhotos,
				uint16(0), // nbVideos,
				uint16(0), // nbPuds,
				uint16(0), // nbCrashLogs,
			},
			c.deprecatedMassStorageContentChanged,
		),
		arcommands.NewD2CCommand(
			12,
			"MassStorageContent",
			[]interface{}{
				uint8(0),  // mass_storage_id,
				uint16(0), // nbPhotos,
				uint16(0), // nbVideos,
				uint16(0), // nbPuds,
				uint16(0), // nbCrashLogs,
				uint16(0), // nbRawPhotos,
			},
			c.massStorageContent,
		),
		arcommands.NewD2CCommand(
			13,
			"MassStorageContentForCurrentRun",
			[]interface{}{
				uint8(0),  // mass_storage_id,
				uint16(0), // nbPhotos,
				uint16(0), // nbVideos,
				uint16(0), // nbRawPhotos,
			},
			c.massStorageContentForCurrentRun,
		),
		arcommands.NewD2CCommand(
			14,
			"VideoRecordingTimestamp",
			[]interface{}{
				uint64(0), // startTimestamp,
				uint64(0), // stopTimestamp,
			},
			c.videoRecordingTimestamp,
		),
	}
}

// allStatesChanged is invoked by the device when all states have been sent.
func (c *commonState) allStatesChanged(args []interface{}) error {
	c.Lock()
	defer c.Unlock()
	c.allStatesSent = ptr.ToBool(true)
	log.Debug("all states have been sent by the device")
	return nil
}

// batteryStateChanged is invoked by the device when the battery state changes.
func (c *commonState) batteryStateChanged(args []interface{}) error {
	c.Lock()
	defer c.Unlock()
	c.batteryPercent = ptr.ToUint8(args[0].(uint8))
	log.WithField(
		"batteryPercent", *c.batteryPercent,
	).Debug("battery state changed")
	return nil
}

// massStorageStateListChanged is invoked by the device when a mass storage
// device is inserted or ejected.
func (c *commonState) massStorageStateListChanged(args []interface{}) error {
	c.Lock()
	defer c.Unlock()
	massStorageID := args[0].(uint8)
	msdIface, ok := c.massStorageDevices[massStorageID]
	var msd *massStorageDevice
	if ok {
		msd = msdIface.(*massStorageDevice)
	} else {
		msd = &massStorageDevice{
			id: massStorageID,
		}
		c.massStorageDevices[massStorageID] = msd
	}
	msd.name = ptr.ToString(args[1].(string))
	log.WithField(
		"id", msd.id,
	).WithField(
		"name", *msd.name,
	).Debug("mass storage state list changed")
	return nil
}

// massStorageInfoStateListChanged is invoked by the devicde when mass storage
// info changes.
func (c *commonState) massStorageInfoStateListChanged(args []interface{}) error {
	c.Lock()
	defer c.Unlock()
	massStorageID := args[0].(uint8)
	msdIface, ok := c.massStorageDevices[massStorageID]
	var msd *massStorageDevice
	if ok {
		msd = msdIface.(*massStorageDevice)
	} else {
		msd = &massStorageDevice{
			id: massStorageID,
		}
		c.massStorageDevices[massStorageID] = msd
	}
	msd.size = ptr.ToUint32(args[1].(uint32))
	msd.usedSize = ptr.ToUint32(args[2].(uint32))
	msd.plugged = ptr.ToBool(args[3].(uint8) == 1)
	msd.full = ptr.ToBool(args[4].(uint8) == 1)
	msd.internal = ptr.ToBool(args[5].(uint8) == 1)
	log.WithField(
		"id", msd.id,
	).WithField(
		"size", *msd.size,
	).WithField(
		"usedSize", *msd.usedSize,
	).WithField(
		"plugged", *msd.plugged,
	).WithField(
		"full", *msd.full,
	).WithField(
		"internal", *msd.internal,
	).Debug("common.massStorageInfoStateListChanged() called")
	return nil
}

// TODO: Implement this
// Title: Date changed
// Description: Date changed.\n Corresponds to the latest date set on the drone.
//   \n\n **Please note that you should not care about this event if you are
//   using the libARController API as this library is handling the connection
//   process for you.**
// Support: drones
// Triggered: by [SetDate](#0-4-1).
// Result:
func (c *commonState) currentDateChanged(args []interface{}) error {
	// date := args[0].(string)
	//   Date with ISO-8601 format
	log.Info("common.currentDateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Time changed
// Description: Time changed.\n Corresponds to the latest time set on the drone.
//   \n\n **Please note that you should not care about this event if you are
//   using the libARController API as this library is handling the connection
//   process for you.**
// Support: drones
// Triggered: by [SetTime](#0-4-2).
// Result:
func (c *commonState) currentTimeChanged(args []interface{}) error {
	// time := args[0].(string)
	//   Time with ISO-8601 format
	log.Info("common.currentTimeChanged() called")
	return nil
}

// // TODO: Implement this
// // Title: Mass storage remaining data list
// // Description: Mass storage remaining data list.
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (c *commonState) massStorageInfoRemainingListChanged(
// 	args []interface{},
// ) error {
// 	// free_space := args[0].(uint32)
// 	//   Mass storage free space in MBytes
// 	// rec_time := args[1].(uint16)
// 	//   Mass storage record time reamining in minute
// 	// photo_remaining := args[2].(uint32)
// 	//   Mass storage photo remaining
// 	log.Info("common.massStorageInfoRemainingListChanged() called")
// 	return nil
// }

// wifiSignalChanged is invoked when the device reports relative wifi signal
// strength at regular intervals.
func (c *commonState) wifiSignalChanged(args []interface{}) error {
	c.Lock()
	defer c.Unlock()
	c.rssi = ptr.ToInt16(args[0].(int16))
	log.WithField(
		"rssi", *c.rssi,
	).Debug("common state wifi signal strength updated")
	return nil
}

// TODO: Implement this
// Title: Sensors state list
// Description: Sensors state list.
// Support: 0901:2.0.3;0902;0905;0906;0907;0909;090a;090c;090e
// Triggered: at connection and when a sensor state changes.
// Result:
func (c *commonState) sensorsStatesListChanged(args []interface{}) error {
	c.Lock()
	defer c.Unlock()
	sensorID := args[0].(int32)
	state := ptr.ToBool(args[1].(uint8) == 1)
	log := log.WithField(
		"state", *state,
	)
	switch sensorID {
	case 0:
		c.imuOK = state
		log = log.WithField("sensor", "IMU")
	case 1:
		c.barometerOK = state
		log = log.WithField("sensor", "barometer")
	case 2:
		c.ultrasoundOK = state
		log = log.WithField("sensor", "ultrasound")
	case 3:
		c.gpsOK = state
		log = log.WithField("sensor", "GPS")
	case 4:
		c.magnetomenterOK = state
		log = log.WithField("sensor", "megnetometer")
	case 5:
		c.verticalCameraOK = state
		log = log.WithField("sensor", "verticalCamera")
	default:
		log.WithField(
			"sensorID", sensorID,
		).Warn("sensor state changed with unknown sensor id")
		return nil
	}
	log.Debug("sensor state changed")
	return nil
}

// // TODO: Implement this
// // Title: Product sub-model
// // Description: Product sub-model.\n This can be used to customize the UI
// //   depending on the product.
// // Support: 0905;0906;0907;0909
// // Triggered: at connection.
// // Result:
// func (c *commonState) productModel(args []interface{}) error {
// 	// model := args[0].(int32)
// 	//   The Model of the product.
// 	//   0: RS_TRAVIS: Travis (RS taxi) model.
// 	//   1: RS_MARS: Mars (RS space) model
// 	//   2: RS_SWAT: SWAT (RS SWAT) model
// 	//   3: RS_MCLANE: Mc Lane (RS police) model
// 	//   4: RS_BLAZE: Blaze (RS fire) model
// 	//   5: RS_ORAK: Orak (RS carbon hydrofoil) model
// 	//   6: RS_NEWZ: New Z (RS wooden hydrofoil) model
// 	//   7: JS_MARSHALL: Marshall (JS fire) model
// 	//   8: JS_DIESEL: Diesel (JS SWAT) model
// 	//   9: JS_BUZZ: Buzz (JS space) model
// 	//   10: JS_MAX: Max (JS F1) model
// 	//   11: JS_JETT: Jett (JS flames) model
// 	//   12: JS_TUKTUK: Tuk-Tuk (JS taxi) model
// 	//   13: SW_BLACK: Swing black model
// 	//   14: SW_WHITE: Swing white model
// 	log.Info("common.productModel() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Country list
// // Description: List of countries known by the drone.
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (c *commonState) countryListKnown(args []interface{}) error {
// 	// listFlags := args[0].(uint8)
// 	//   List entry attribute Bitfield. 0x01: First: indicate it&#39;s the first
// 	//   element of the list. 0x02: Last: indicate it&#39;s the last element of
// 	//   the list. 0x04: Empty: indicate the list is empty (implies First/Last).
// 	//   All other arguments should be ignored.
// 	// countryCodes := args[1].(string)
// 	//   Following of country code with ISO 3166 format, separated by &#34;;&#34;.
// 	//   Be careful of the command size allowed by the network used. If necessary,
// 	//   split the list in several commands.
// 	log.Info("common.countryListKnown() called")
// 	return nil
// }

// deprecatedMassStorageContentChanged is deprecated in favor of
// massStorageContent, but since we can still see this command being invoked,
// we'll implement the command to avoid a warning, but the implementation will
// remain a no-op unless / until such time that it becomes clear that older
// versions of the firmware might require us to support both commands.
func (c *commonState) deprecatedMassStorageContentChanged(
	args []interface{},
) error {
	// mass_storage_id := args[0].(uint8)
	//   Mass storage id (unique)
	// nbPhotos := args[1].(uint16)
	//   Number of photos (does not include raw photos)
	// nbVideos := args[2].(uint16)
	//   Number of videos
	// nbPuds := args[3].(uint16)
	//   Number of puds
	// nbCrashLogs := args[4].(uint16)
	//   Number of crash logs
	log.Debug("mass storage content changed-- this is a no-op")
	return nil
}

// massStorageContent is invoked when the device reports that the content of the
// mass storage has changed.
func (c *commonState) massStorageContent(args []interface{}) error {
	c.Lock()
	defer c.Unlock()
	massStorageID := args[0].(uint8)
	msdIface, ok := c.massStorageDevices[massStorageID]
	var msd *massStorageDevice
	if ok {
		msd = msdIface.(*massStorageDevice)
	} else {
		msd = &massStorageDevice{
			id: massStorageID,
		}
		c.massStorageDevices[massStorageID] = msd
	}
	msd.photoCount = ptr.ToUint16(args[1].(uint16))
	msd.videoCount = ptr.ToUint16(args[2].(uint16))
	msd.pudCount = ptr.ToUint16(args[3].(uint16))
	msd.crashLogCount = ptr.ToUint16(args[4].(uint16))
	msd.rawPhotoCount = ptr.ToUint16(args[5].(uint16))
	log.WithField(
		"id", msd.id,
	).WithField(
		"photoCount", *msd.photoCount,
	).WithField(
		"videoCount", *msd.videoCount,
	).WithField(
		"pudCount", *msd.pudCount,
	).WithField(
		"crashLogCount", *msd.crashLogCount,
	).WithField(
		"rawPhotoCount", *msd.rawPhotoCount,
	).Debug("mass storage content changed")
	return nil
}

// massStorageContentForCurrentRun is invoked when the device reports that the
// content of the mass storage has changed and the content is related to the
// current run.
func (c *commonState) massStorageContentForCurrentRun(
	args []interface{},
) error {
	c.Lock()
	defer c.Unlock()
	massStorageID := args[0].(uint8)
	msdIface, ok := c.massStorageDevices[massStorageID]
	var msd *massStorageDevice
	if ok {
		msd = msdIface.(*massStorageDevice)
	} else {
		msd = &massStorageDevice{
			id: massStorageID,
		}
		c.massStorageDevices[massStorageID] = msd
	}
	msd.currentRunPhotoCount = ptr.ToUint16(args[1].(uint16))
	msd.currentRunVideoCount = ptr.ToUint16(args[2].(uint16))
	msd.currentRunRawPhotoCount = ptr.ToUint16(args[3].(uint16))
	log.WithField(
		"id", msd.id,
	).WithField(
		"currentRunPhotoCount", *msd.currentRunPhotoCount,
	).WithField(
		"currentRunVideoCount", *msd.currentRunVideoCount,
	).WithField(
		"currentRunRawPhotoCount", *msd.currentRunRawPhotoCount,
	).Debug("mass storage content for current run changed")
	return nil
}

// videoRecordingTimestamp is invoked by the device on video recording start and
// video recording stop or after that the date/time of the drone changed.
func (c *commonState) videoRecordingTimestamp(args []interface{}) error {
	// startTimestamp := args[0].(uint64)
	//   Timestamp in milliseconds since 00:00:00 UTC on 1 January 1970.
	// stopTimestamp := args[1].(uint64)
	//   Timestamp in milliseconds since 00:00:00 UTC on 1 January 1970. 0 mean
	//   that video is still recording.
	log.Info("common.videoRecordingTimestamp() called")
	return nil
}

func (c *commonState) RSSI() (int16, bool) {
	if c.rssi == nil {
		return 0, false
	}
	return *c.rssi, true
}

func (c *commonState) BatteryPercent() (uint8, bool) {
	if c.batteryPercent == nil {
		return 0, false
	}
	return *c.batteryPercent, true
}

func (c *commonState) IMUOK() (bool, bool) {
	if c.imuOK == nil {
		return false, false
	}
	return *c.imuOK, true
}

func (c *commonState) BarometerOK() (bool, bool) {
	if c.barometerOK == nil {
		return false, false
	}
	return *c.barometerOK, true
}

func (c *commonState) UltrasoundOK() (bool, bool) {
	if c.ultrasoundOK == nil {
		return false, false
	}
	return *c.ultrasoundOK, true
}

func (c *commonState) GPSOK() (bool, bool) {
	if c.gpsOK == nil {
		return false, false
	}
	return *c.gpsOK, true
}

func (c *commonState) MagnetometerOK() (bool, bool) {
	if c.magnetomenterOK == nil {
		return false, false
	}
	return *c.magnetomenterOK, true
}

func (c *commonState) VerticalCameraOK() (bool, bool) {
	if c.verticalCameraOK == nil {
		return false, false
	}
	return *c.verticalCameraOK, true
}

func (c *commonState) AllStatesSent() (bool, bool) {
	if c.allStatesSent == nil {
		return false, false
	}
	return *c.allStatesSent, true
}

func (c *commonState) MassStorageDevices() map[uint8]MassStorageDevice {
	return c.massStorageDevices
}
