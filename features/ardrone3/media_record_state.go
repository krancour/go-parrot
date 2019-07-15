package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// State of media recording

const (
	PictureStateReady       int32 = 0
	PictureStateBusy        int32 = 1
	PictureStateUnavailable int32 = 2

	PictureStateReasonOK         int32 = 0
	PictureStateReasonUnknown    int32 = 1 // Unknown, generic error
	PictureStateReasonCameraKO   int32 = 2 // Still camera is out of order
	PictureStateReasonMemoryFull int32 = 3 // Memory is full; cannot save more pictures
	PictureStateReasonLowBattery int32 = 4 // Battery is too low to take a picture

	VideoStateStopped      int32 = 0
	VideoStateStarted      int32 = 1
	VideoStateNotAvailable int32 = 2

	VideoStateReasonOK         int32 = 0
	VideoStateReasonUnknown    int32 = 1 // Unknown, generic error
	VideoStateReasonCameraKO   int32 = 2 // Video camera is out of order
	VideoStateReasonMemoryFull int32 = 3 // Memory is full; cannot save more video
	VideoStateReasonLowBattery int32 = 4 // Battery is too low to start/continue recording
)

// MediaRecordState ...
// TODO: Document this
type MediaRecordState interface {
	lock.ReadLockable
	// PictureState returns the still camera state. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	PictureState() (int32, bool)
	// PictureStateReason returns the reason for the still camera state. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	PictureStateReason() (int32, bool)
	// VideoState returns the video camera state. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	VideoState() (int32, bool)
	// VideoStateReason return the reason for the video camera state. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	VideoStateReason() (int32, bool)
}

type mediaRecordState struct {
	sync.RWMutex
	// pictureState represents the still camera state
	pictureState *int32
	// pictureStateReason represents the reason for the still camera state
	pictureStateReason *int32
	// videoState represents the video camera state
	videoState *int32
	// videoStateReason represents the reason for the video camera state
	videoStateReason *int32
}

func newMediaRecordState() *mediaRecordState {
	return &mediaRecordState{}
}

func (m *mediaRecordState) ID() uint8 {
	return 8
}

func (m *mediaRecordState) Name() string {
	return "MediaRecordState"
}

func (m *mediaRecordState) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			2,
			"PictureStateChangedV2",
			[]interface{}{
				int32(0), // state,
				int32(0), // error,
			},
			m.pictureStateChangedV2,
			log,
		),
		arcommands.NewD2CCommand(
			3,
			"VideoStateChangedV2",
			[]interface{}{
				int32(0), // state,
				int32(0), // error,
			},
			m.videoStateChangedV2,
			log,
		),
	}
}

// pictureStateChangedV2 is invoked by the device when the still camera's
// availability changes.
func (m *mediaRecordState) pictureStateChangedV2(
	args []interface{},
	log *log.Entry,
) error {
	m.Lock()
	defer m.Unlock()
	m.pictureState = ptr.ToInt32(args[0].(int32))
	m.pictureStateReason = ptr.ToInt32(args[1].(int32))
	log.WithField(
		"state", *m.pictureState,
	).WithField(
		"error", *m.pictureStateReason,
	).Debug("picture state changed")
	return nil
}

// videoStateChangedV2 is invoked by the device when there is a change in
// video camera state or availability.
func (m *mediaRecordState) videoStateChangedV2(
	args []interface{},
	log *log.Entry,
) error {
	m.Lock()
	defer m.Unlock()
	m.videoState = ptr.ToInt32(args[0].(int32))
	m.videoStateReason = ptr.ToInt32(args[1].(int32))
	log.WithField(
		"state", *m.videoState,
	).WithField(
		"error", *m.videoStateReason,
	).Debug("video state changed")
	return nil
}

func (m *mediaRecordState) PictureState() (int32, bool) {
	if m.pictureState == nil {
		return 0, false
	}
	return *m.pictureState, true
}

func (m *mediaRecordState) PictureStateReason() (int32, bool) {
	if m.pictureStateReason == nil {
		return 0, false
	}
	return *m.pictureStateReason, true
}

func (m *mediaRecordState) VideoState() (int32, bool) {
	if m.videoState == nil {
		return 0, false
	}
	return *m.videoState, true
}

func (m *mediaRecordState) VideoStateReason() (int32, bool) {
	if m.videoStateReason == nil {
		return 0, false
	}
	return *m.videoStateReason, true
}
