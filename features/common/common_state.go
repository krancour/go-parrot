package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Common state from product

// CommonState ...
// TODO: Document this
type CommonState interface {
	// RLock blocks until a read lock is obtained. This permits callers to procede
	// with querying any or all attributes of the common state without worry
	// that some attributes will be overwritten as others are read. i.e. It
	// permits the possibility of taking an atomic snapshop of common state.
	// Note that use of this function is not obligatory for applications that do
	// not require such guarantees. Callers MUST call RUnlock() or else piloting
	// state will never resume updating.
	RLock()
	// RUnlock releases a read lock on the common state. See RLock().
	RUnlock()
	// RSSI returns the relative signal stength between the client and the device
	// in dbm. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	RSSI() (int16, bool)
}

type commonState struct {
	// TODO: Is this right? I thought RSSI is a relative measure, while dbm
	// would seem to indicate an absolute measure.
	// rssi is the relative signal stength between the client and the device
	// in dbm
	rssi *int16
	lock sync.RWMutex
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
		// arcommands.NewD2CCommand(
		// 	11,
		// 	"DeprecatedMassStorageContentChanged",
		// 	[]interface{}{
		// 		uint8(0),  // mass_storage_id,
		// 		uint16(0), // nbPhotos,
		// 		uint16(0), // nbVideos,
		// 		uint16(0), // nbPuds,
		// 		uint16(0), // nbCrashLogs,
		// 	},
		// 	c.deprecatedMassStorageContentChanged,
		// ),
		// arcommands.NewD2CCommand(
		// 	12,
		// 	"MassStorageContent",
		// 	[]interface{}{
		// 		uint8(0),  // mass_storage_id,
		// 		uint16(0), // nbPhotos,
		// 		uint16(0), // nbVideos,
		// 		uint16(0), // nbPuds,
		// 		uint16(0), // nbCrashLogs,
		// 		uint16(0), // nbRawPhotos,
		// 	},
		// 	c.massStorageContent,
		// ),
		// arcommands.NewD2CCommand(
		// 	13,
		// 	"MassStorageContentForCurrentRun",
		// 	[]interface{}{
		// 		uint8(0),  // mass_storage_id,
		// 		uint16(0), // nbPhotos,
		// 		uint16(0), // nbVideos,
		// 		uint16(0), // nbRawPhotos,
		// 	},
		// 	c.massStorageContentForCurrentRun,
		// ),
		// arcommands.NewD2CCommand(
		// 	14,
		// 	"VideoRecordingTimestamp",
		// 	[]interface{}{
		// 		uint64(0), // startTimestamp,
		// 		uint64(0), // stopTimestamp,
		// 	},
		// 	c.videoRecordingTimestamp,
		// ),
	}
}

// TODO: Implement this
// Title: All states have been sent
// Description: All states have been sent.\n\n **Please note that you should not
//   care about this event if you are using the libARController API as this
//   library is handling the connection process for you.**
// Support: drones
// Triggered: when all states values have been sent.
// Result:
func (c *commonState) allStatesChanged(args []interface{}) error {
	log.Info("common.allStatesChanged() called")
	return nil
}

// TODO: Implement this
// Title: Battery state
// Description: Battery state.
// Support: drones
// Triggered: when the battery level changes.
// Result:
func (c *commonState) batteryStateChanged(args []interface{}) error {
	// percent := args[0].(uint8)
	//   Battery percentage
	log.Info("common.batteryStateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Mass storage state list
// Description: Mass storage state list.
// Support: drones
// Triggered: when a mass storage is inserted or ejected.
// Result:
func (c *commonState) massStorageStateListChanged(args []interface{}) error {
	// mass_storage_id := args[0].(uint8)
	//   Mass storage id (unique)
	// name := args[1].(string)
	//   Mass storage name
	log.Info("common.massStorageStateListChanged() called")
	return nil
}

// TODO: Implement this
// Title: Mass storage info state list
// Description: Mass storage info state list.
// Support: drones
// Triggered: when a mass storage info changes.
// Result:
func (c *commonState) massStorageInfoStateListChanged(args []interface{}) error {
	// mass_storage_id := args[0].(uint8)
	//   Mass storage state id (unique)
	// size := args[1].(uint32)
	//   Mass storage size in MBytes
	// used_size := args[2].(uint32)
	//   Mass storage used size in MBytes
	// plugged := args[3].(uint8)
	//   Mass storage plugged (1 if mass storage is plugged, otherwise 0)
	// full := args[4].(uint8)
	//   Mass storage full information state (1 if mass storage full, 0
	//   otherwise).
	// internal := args[5].(uint8)
	//   Mass storage internal type state (1 if mass storage is internal, 0
	//   otherwise)
	log.Info("common.massStorageInfoStateListChanged() called")
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
	c.lock.Lock()
	defer c.lock.Unlock()
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
	// sensorName := args[0].(int32)
	//   Sensor name
	//   0: IMU: Inertial Measurement Unit sensor
	//   1: barometer: Barometer sensor
	//   2: ultrasound: Ultrasonic sensor
	//   3: GPS: GPS sensor
	//   4: magnetometer: Magnetometer sensor
	//   5: vertical_camera: Vertical Camera sensor
	// sensorState := args[1].(uint8)
	//   Sensor state (1 if the sensor is OK, 0 if the sensor is NOT OK)
	log.Info("common.sensorsStatesListChanged() called")
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

// // TODO: Implement this
// // Title: Mass storage content changed
// // Description: Mass storage content changed.
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (c *commonState) deprecatedMassStorageContentChanged(
// 	args []interface{},
// ) error {
// 	// mass_storage_id := args[0].(uint8)
// 	//   Mass storage id (unique)
// 	// nbPhotos := args[1].(uint16)
// 	//   Number of photos (does not include raw photos)
// 	// nbVideos := args[2].(uint16)
// 	//   Number of videos
// 	// nbPuds := args[3].(uint16)
// 	//   Number of puds
// 	// nbCrashLogs := args[4].(uint16)
// 	//   Number of crash logs
// 	log.Info("common.deprecatedMassStorageContentChanged() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Mass storage content
// // Description: Mass storage content.
// // Support: 090c:4.0.0;090e:4.0.0
// // Triggered: when the content of the mass storage changes.
// // Result:
// func (c *commonState) massStorageContent(args []interface{}) error {
// 	// mass_storage_id := args[0].(uint8)
// 	//   Mass storage id (unique)
// 	// nbPhotos := args[1].(uint16)
// 	//   Number of photos (does not include raw photos)
// 	// nbVideos := args[2].(uint16)
// 	//   Number of videos
// 	// nbPuds := args[3].(uint16)
// 	//   Number of puds
// 	// nbCrashLogs := args[4].(uint16)
// 	//   Number of crash logs
// 	// nbRawPhotos := args[5].(uint16)
// 	//   Number of raw photos
// 	log.Info("common.massStorageContent() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Mass storage content for current run
// // Description: Mass storage content for current run.\n Only counts the files
// //   related to the current run (see [RunId](#0-30-0))
// // Support: 090c:4.0.0;090e:4.0.0
// // Triggered: when the content of the mass storage changes and this content is
// //   related to the current run.
// // Result:
// func (c *commonState) massStorageContentForCurrentRun(args []interface{}) error {
// 	// mass_storage_id := args[0].(uint8)
// 	//   Mass storage id (unique)
// 	// nbPhotos := args[1].(uint16)
// 	//   Number of photos (does not include raw photos)
// 	// nbVideos := args[2].(uint16)
// 	//   Number of videos
// 	// nbRawPhotos := args[3].(uint16)
// 	//   Number of raw photos
// 	log.Info("common.massStorageContentForCurrentRun() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Video recording timestamp
// // Description: Current or last video recording timestamp.\n Timestamp in
// //   milliseconds since 00:00:00 UTC on 1 January 1970.\n **Please note that
// //   values don&#39;t persist after drone reboot**
// // Support:
// // Triggered: on video recording start and video recording stop or \n after that
// //   the date/time of the drone changed.
// // Result:
// func (c *commonState) videoRecordingTimestamp(args []interface{}) error {
// 	// startTimestamp := args[0].(uint64)
// 	//   Timestamp in milliseconds since 00:00:00 UTC on 1 January 1970.
// 	// stopTimestamp := args[1].(uint64)
// 	//   Timestamp in milliseconds since 00:00:00 UTC on 1 January 1970. 0 mean
// 	//   that video is still recording.
// 	log.Info("common.videoRecordingTimestamp() called")
// 	return nil
// }

func (c *commonState) RLock() {
	c.lock.RLock()
}

func (c *commonState) RUnlock() {
	c.lock.RUnlock()
}

func (c *commonState) RSSI() (int16, bool) {
	if c.rssi == nil {
		return 0, false
	}
	return *c.rssi, true
}
