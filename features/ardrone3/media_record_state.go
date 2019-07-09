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
	}
}

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
