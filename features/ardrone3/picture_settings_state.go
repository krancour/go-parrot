package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Photo settings state from product

const (
	PictureFormatRaw         int32 = 0
	PictureFormatJPEG        int32 = 1
	PictureFormatSnapshot    int32 = 2
	PictureFormatJPEGFisheye int32 = 3
)

// PictureSettingsState ...
// TODO: Document this
type PictureSettingsState interface {
	// RLock blocks until a read lock is obtained. This permits callers to procede
	// with querying any or all attributes of the picture settings state without
	// worry that some attributes will be overwritten as others are read. i.e. It
	// permits the possibility of taking an atomic snapshop of picture settings
	// state. Note that use of this function is not obligatory for applications
	// that do not require such guarantees. Callers MUST call RUnlock() or else
	// picture settings state will never resume updating.
	RLock()
	// RUnlock releases a read lock on the GPS state. See RLock().
	RUnlock()
	// Format returns the picture format. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	Format() (int32, bool)
}

type pictureSettingsState struct {
	// format represents the picture format
	format *int32
	lock   sync.RWMutex
}

func (p *pictureSettingsState) ID() uint8 {
	return 20
}

func (p *pictureSettingsState) Name() string {
	return "PictureSettingsState"
}

func (p *pictureSettingsState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"PictureFormatChanged",
			[]interface{}{
				int32(0), // type,
			},
			p.pictureFormatChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"AutoWhiteBalanceChanged",
			[]interface{}{
				int32(0), // type,
			},
			p.autoWhiteBalanceChanged,
		),
		arcommands.NewD2CCommand(
			2,
			"ExpositionChanged",
			[]interface{}{
				float32(0), // value,
				float32(0), // min,
				float32(0), // max,
			},
			p.expositionChanged,
		),
		arcommands.NewD2CCommand(
			3,
			"SaturationChanged",
			[]interface{}{
				float32(0), // value,
				float32(0), // min,
				float32(0), // max,
			},
			p.saturationChanged,
		),
		arcommands.NewD2CCommand(
			4,
			"TimelapseChanged",
			[]interface{}{
				uint8(0),   // enabled,
				float32(0), // interval,
				float32(0), // minInterval,
				float32(0), // maxInterval,
			},
			p.timelapseChanged,
		),
		arcommands.NewD2CCommand(
			5,
			"VideoAutorecordChanged",
			[]interface{}{
				uint8(0), // enabled,
				uint8(0), // mass_storage_id,
			},
			p.videoAutorecordChanged,
		),
		arcommands.NewD2CCommand(
			6,
			"VideoStabilizationModeChanged",
			[]interface{}{
				int32(0), // mode,
			},
			p.videoStabilizationModeChanged,
		),
		arcommands.NewD2CCommand(
			7,
			"VideoRecordingModeChanged",
			[]interface{}{
				int32(0), // mode,
			},
			p.videoRecordingModeChanged,
		),
		arcommands.NewD2CCommand(
			8,
			"VideoFramerateChanged",
			[]interface{}{
				int32(0), // framerate,
			},
			p.videoFramerateChanged,
		),
		arcommands.NewD2CCommand(
			9,
			"VideoResolutionsChanged",
			[]interface{}{
				int32(0), // type,
			},
			p.videoResolutionsChanged,
		),
	}
}

// pictureFormatChanged is invoked by the device when the picture format is
// changed.
func (p *pictureSettingsState) pictureFormatChanged(args []interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.format = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"type", *p.format,
	).Debug("picture format changed")
	return nil
}

// TODO: Implement this
// Title: White balance mode
// Description: White balance mode.
// Support: 0901;090c;090e
// Triggered: by [SetWhiteBalanceMode](#1-19-1).
// Result:
func (p *pictureSettingsState) autoWhiteBalanceChanged(
	args []interface{},
) error {
	// type := args[0].(int32)
	//   The type auto white balance
	//   0: auto: Auto guess of best white balance params
	//   1: tungsten: Tungsten white balance
	//   2: daylight: Daylight white balance
	//   3: cloudy: Cloudy white balance
	//   4: cool_white: White balance for a flash
	log.Info("ardrone3.autoWhiteBalanceChanged() called")
	return nil
}

// TODO: Implement this
// Title: Image exposure
// Description: Image exposure.
// Support: 0901;090c;090e
// Triggered: by [SetImageExposure](#1-19-2).
// Result:
func (p *pictureSettingsState) expositionChanged(args []interface{}) error {
	// value := args[0].(float32)
	//   Exposure value
	// min := args[1].(float32)
	//   Min exposure value
	// max := args[2].(float32)
	//   Max exposure value
	log.Info("ardrone3.expositionChanged() called")
	return nil
}

// TODO: Implement this
// Title: Image saturation
// Description: Image saturation.
// Support: 0901;090c;090e
// Triggered: by [SetImageSaturation](#1-19-3).
// Result:
func (p *pictureSettingsState) saturationChanged(args []interface{}) error {
	// value := args[0].(float32)
	//   Saturation value
	// min := args[1].(float32)
	//   Min saturation value
	// max := args[2].(float32)
	//   Max saturation value
	log.Info("ardrone3.saturationChanged() called")
	return nil
}

// TODO: Implement this
// Title: Timelapse mode
// Description: Timelapse mode.
// Support: 0901;090c;090e
// Triggered: by [SetTimelapseMode](#1-19-4).
// Result:
func (p *pictureSettingsState) timelapseChanged(args []interface{}) error {
	// enabled := args[0].(uint8)
	//   1 if timelapse is enabled, 0 otherwise
	// interval := args[1].(float32)
	//   interval in seconds for taking pictures
	// minInterval := args[2].(float32)
	//   Minimal interval for taking pictures
	// maxInterval := args[3].(float32)
	//   Maximal interval for taking pictures
	log.Info("ardrone3.timelapseChanged() called")
	return nil
}

// TODO: Implement this
// Title: Video Autorecord mode
// Description: Video Autorecord mode.
// Support: 0901;090c;090e
// Triggered: by [SetVideoAutorecordMode](#1-19-5).
// Result:
func (p *pictureSettingsState) videoAutorecordChanged(args []interface{}) error {
	// enabled := args[0].(uint8)
	//   1 if video autorecord is enabled, 0 otherwise
	// mass_storage_id := args[1].(uint8)
	//   Mass storage id for the taken video
	log.Info("ardrone3.videoAutorecordChanged() called")
	return nil
}

// TODO: Implement this
// Title: Video stabilization mode
// Description: Video stabilization mode.
// Support: 0901:3.4.0;090c:3.4.0;090e
// Triggered: by [SetVideoStabilizationMode](#1-19-6).
// Result:
func (p *pictureSettingsState) videoStabilizationModeChanged(
	args []interface{},
) error {
	// mode := args[0].(int32)
	//   Video stabilization mode
	//   0: roll_pitch: Video flat on roll and pitch
	//   1: pitch: Video flat on pitch only
	//   2: roll: Video flat on roll only
	//   3: none: Video follows drone angles
	log.Info("ardrone3.videoStabilizationModeChanged() called")
	return nil
}

// TODO: Implement this
// Title: Video recording mode
// Description: Video recording mode.
// Support: 0901:3.4.0;090c:3.4.0;090e
// Triggered: by [SetVideoRecordingMode](#1-19-7).
// Result:
func (p *pictureSettingsState) videoRecordingModeChanged(
	args []interface{},
) error {
	// mode := args[0].(int32)
	//   Video recording mode
	//   0: quality: Maximize recording quality.
	//   1: time: Maximize recording time.
	log.Info("ardrone3.videoRecordingModeChanged() called")
	return nil
}

// TODO: Implement this
// Title: Video framerate
// Description: Video framerate.
// Support: 0901:3.4.0;090c:3.4.0;090e
// Triggered: by [SetVideoFramerateMode](#1-19-8).
// Result:
func (p *pictureSettingsState) videoFramerateChanged(args []interface{}) error {
	// framerate := args[0].(int32)
	//   Video framerate
	//   0: 24_FPS: 23.976 frames per second.
	//   1: 25_FPS: 25 frames per second.
	//   2: 30_FPS: 29.97 frames per second.
	log.Info("ardrone3.videoFramerateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Video resolutions
// Description: Video resolutions.\n This event informs about the recording AND
//   streaming resolutions.
// Support: 0901:3.4.0;090c:3.4.0;090e
// Triggered: by [SetVideResolutions](#1-19-9).
// Result:
func (p *pictureSettingsState) videoResolutionsChanged(
	args []interface{},
) error {
	// type := args[0].(int32)
	//   Video resolution type.
	//   0: rec1080_stream480: 1080p recording, 480p streaming.
	//   1: rec720_stream720: 720p recording, 720p streaming.
	log.Info("ardrone3.videoResolutionsChanged() called")
	return nil
}

func (p *pictureSettingsState) RLock() {
	p.lock.RLock()
}

func (p *pictureSettingsState) RUnlock() {
	p.lock.RUnlock()
}

func (p *pictureSettingsState) Format() (int32, bool) {
	if p.format == nil {
		return 0, false
	}
	return *p.format, true
}
