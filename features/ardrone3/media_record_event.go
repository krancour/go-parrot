package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Events of media recording

// MediaRecordEvent ...
// TODO: Document this
type MediaRecordEvent interface {
	lock.ReadLockable
}

type mediaRecordEvent struct {
	sync.RWMutex
}

func (m *mediaRecordEvent) ID() uint8 {
	return 3
}

func (m *mediaRecordEvent) Name() string {
	return "MediaRecordEvent"
}

func (m *mediaRecordEvent) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"PictureEventChanged",
			[]interface{}{
				int32(0), // event,
				int32(0), // error,
			},
			m.pictureEventChanged,
			log,
		),
		arcommands.NewD2CCommand(
			1,
			"VideoEventChanged",
			[]interface{}{
				int32(0), // event,
				int32(0), // error,
			},
			m.videoEventChanged,
			log,
		),
	}
}

// TODO: Implement this
// Title: Picture taken
// Description: Picture taken.\n\n **This event is a notification, you can&#39;t
//   retrieve it in the cache of the device controller.**
// Support: 0901:2.0.1;090c;090e
// Triggered: after a [TakePicture](#1-7-2), when the picture has been taken
//   (or it has failed).
// Result:
func (m *mediaRecordEvent) pictureEventChanged(
	args []interface{},
	log *log.Entry,
) error {
	// event := args[0].(int32)
	//   Last event of picture recording
	//   0: taken: Picture taken and saved
	//   1: failed: Picture failed
	// error := args[1].(int32)
	//   Error to explain the event
	//   0: ok: No Error
	//   1: unknown: Unknown generic error ; only when state is failed
	//   2: busy: Picture recording is busy ; only when state is failed
	//   3: notAvailable: Picture recording not available ; only when state is
	//      failed
	//   4: memoryFull: Memory full ; only when state is failed
	//   5: lowBattery: Battery is too low to record.
	log.Warn("command not implemented")
	return nil
}

// TODO: Implement this
// Title: Video record notification
// Description: Video record notification.\n\n **This event is a notification,
//   you can&#39;t retrieve it in the cache of the device controller.**
// Support: 0901:2.0.1;090c;090e
// Triggered: by [RecordVideo](#1-7-3) or a change in the video state.
// Result:
func (m *mediaRecordEvent) videoEventChanged(
	args []interface{},
	log *log.Entry,
) error {
	// event := args[0].(int32)
	//   Event of video recording
	//   0: start: Video start
	//   1: stop: Video stop and saved
	//   2: failed: Video failed
	// error := args[1].(int32)
	//   Error to explain the event
	//   0: ok: No Error
	//   1: unknown: Unknown generic error ; only when state is failed
	//   2: busy: Video recording is busy ; only when state is failed
	//   3: notAvailable: Video recording not available ; only when state is
	//      failed
	//   4: memoryFull: Memory full
	//   5: lowBattery: Battery is too low to record.
	//   6: autoStopped: Video was auto stopped
	log.Warn("command not implemented")
	return nil
}
