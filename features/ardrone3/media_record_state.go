package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// State of media recording

// MediaRecordState ...
// TODO: Document this
type MediaRecordState interface {
	lock.ReadLockable
}

type mediaRecordState struct {
	sync.RWMutex
}

func (m *mediaRecordState) ID() uint8 {
	return 8
}

func (m *mediaRecordState) Name() string {
	return "MediaRecordState"
}

func (m *mediaRecordState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		// arcommands.NewD2CCommand(
		// 	0,
		// 	"PictureStateChanged",
		// 	[]interface{}{
		// 		uint8(0), // state,
		// 		uint8(0), // mass_storage_id,
		// 	},
		// 	m.pictureStateChanged,
		// ),
		// arcommands.NewD2CCommand(
		// 	1,
		// 	"VideoStateChanged",
		// 	[]interface{}{
		// 		int32(0), // state,
		// 		uint8(0), // mass_storage_id,
		// 	},
		// 	m.videoStateChanged,
		// ),
		arcommands.NewD2CCommand(
			2,
			"PictureStateChangedV2",
			[]interface{}{
				int32(0), // state,
				int32(0), // error,
			},
			m.pictureStateChangedV2,
		),
		arcommands.NewD2CCommand(
			3,
			"VideoStateChangedV2",
			[]interface{}{
				int32(0), // state,
				int32(0), // error,
			},
			m.videoStateChangedV2,
		),
		// arcommands.NewD2CCommand(
		// 	4,
		// 	"VideoResolutionState",
		// 	[]interface{}{
		// 		int32(0), // streaming,
		// 		int32(0), // recording,
		// 	},
		// 	m.videoResolutionState,
		// ),
	}
}

// // TODO: Implement this
// // Title: Picture state
// // Description: Picture state.
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (m *mediaRecordState) pictureStateChanged(args []interface{}) error {
// 	// state := args[0].(uint8)
// 	//   1 if picture has been taken, 0 otherwise
// 	// mass_storage_id := args[1].(uint8)
// 	//   Mass storage id where the picture was recorded
// 	log.Info("ardrone3.pictureStateChanged() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Video record state
// // Description: Picture record state.
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (m *mediaRecordState) videoStateChanged(args []interface{}) error {
// 	// state := args[0].(int32)
// 	//   State of video
// 	//   0: stopped: Video was stopped
// 	//   1: started: Video was started
// 	//   2: failed: Video was failed
// 	//   3: autostopped: Video was auto stopped
// 	// mass_storage_id := args[1].(uint8)
// 	//   Mass storage id where the video was recorded
// 	log.Info("ardrone3.videoStateChanged() called")
// 	return nil
// }

// TODO: Implement this
// Title: Picture state
// Description: Picture state.
// Support: 0901:2.0.1;090c;090e
// Triggered: by [TakePicture](#1-7-2) or by a change in the picture state
// Result:
func (m *mediaRecordState) pictureStateChangedV2(args []interface{}) error {
	// state := args[0].(int32)
	//   State of device picture recording
	//   0: ready: The picture recording is ready
	//   1: busy: The picture recording is busy
	//   2: notAvailable: The picture recording is not available
	// error := args[1].(int32)
	//   Error to explain the state
	//   0: ok: No Error
	//   1: unknown: Unknown generic error
	//   2: camera_ko: Picture camera is out of order
	//   3: memoryFull: Memory full ; cannot save one additional picture
	//   4: lowBattery: Battery is too low to start/keep recording.
	log.Info("ardrone3.pictureStateChangedV2() called")
	return nil
}

// TODO: Implement this
// Title: Video record state
// Description: Video record state.
// Support: 0901:2.0.1;090c;090e
// Triggered: by [RecordVideo](#1-7-3) or by a change in the video state
// Result:
func (m *mediaRecordState) videoStateChangedV2(args []interface{}) error {
	// state := args[0].(int32)
	//   State of device video recording
	//   0: stopped: Video is stopped
	//   1: started: Video is started
	//   2: notAvailable: The video recording is not available
	// error := args[1].(int32)
	//   Error to explain the state
	//   0: ok: No Error
	//   1: unknown: Unknown generic error
	//   2: camera_ko: Video camera is out of order
	//   3: memoryFull: Memory full ; cannot save one additional video
	//   4: lowBattery: Battery is too low to start/keep recording.
	log.Info("ardrone3.videoStateChangedV2() called")
	return nil
}

// // TODO: Implement this
// // Title: Video resolution
// // Description: Video resolution.\n Informs about streaming and recording video
// //   resolutions.\n Note that this is only an indication about what the
// //   resolution should be. To know the real resolution, you should get it from
// //   the frame.
// // Support: none
// // Triggered: when the resolution changes.
// // Result:
// // WARNING: Deprecated
// func (m *mediaRecordState) videoResolutionState(args []interface{}) error {
// 	// streaming := args[0].(int32)
// 	//   Streaming resolution
// 	//   0: res360p: 360p resolution.
// 	//   1: res480p: 480p resolution.
// 	//   2: res720p: 720p resolution.
// 	//   3: res1080p: 1080p resolution.
// 	// recording := args[1].(int32)
// 	//   Recording resolution
// 	//   0: res360p: 360p resolution.
// 	//   1: res480p: 480p resolution.
// 	//   2: res720p: 720p resolution.
// 	//   3: res1080p: 1080p resolution.
// 	log.Info("ardrone3.videoResolutionState() called")
// 	return nil
// }
