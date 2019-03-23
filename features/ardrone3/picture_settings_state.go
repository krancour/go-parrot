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
	// Saturation returns the current image saturation. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	Saturation() (float32, bool)
	// MinSaturation returns the minimum image saturation. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	MinSaturation() (float32, bool)
	// MaxSaturation returns the maximum image saturation. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	MaxSaturation() (float32, bool)
	// Exposure returns the current exposure. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	Exposure() (float32, bool)
	// MinExposure returns the minimum exposure. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	MinExposure() (float32, bool)
	// MaxExposure returns the maximum exposure. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	MaxExposure() (float32, bool)
	// TimeLapseEnabled returns an indicator of whether time lapse image capture
	// is enabled.
	TimeLapseEnabled() (bool, bool)
	// TimeLapseInterval returns the current time lapse interval. A boolean value
	// is also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	TimeLapseInterval() (float32, bool)
	// MinTimeLapseInterval returns the minimum time lapse interval. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	MinTimeLapseInterval() (float32, bool)
	// MaxTimeLapseInterval returns the maximum time lapse interval. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	MaxTimeLapseInterval() (float32, bool)
}

type pictureSettingsState struct {
	// format represents the picture format
	format *int32
	// saturation is the current image saturation
	saturation *float32
	// minSaturation is the minimum image saturation
	minSaturation *float32
	// maxSaturation is the maximum image saturation
	maxSaturation *float32
	// exposure is the current exposure
	exposure *float32
	// minExposure is the minimum exposure
	minExposure *float32
	// maxExposure is the maximum exposure
	maxExposure *float32
	// timeLapseEnabled indicates whether time lapse image capture is enabled
	timeLapseEnabled *bool
	// timeLapseInterval is the current time lapse interval in seconds
	timeLapseInterval *float32
	// minTimeLapseInterval is the minimum time lapse interval in seconds
	minTimeLapseInterval *float32
	// maxTimeLapseInterval is the maximum time lapse interval in seconds
	maxTimeLapseInterval *float32
	lock                 sync.RWMutex
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

// expositionChanged is invoked by the device when exposure is changed.
func (p *pictureSettingsState) expositionChanged(args []interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.exposure = ptr.ToFloat32(args[0].(float32))
	p.minExposure = ptr.ToFloat32(args[1].(float32))
	p.maxExposure = ptr.ToFloat32(args[2].(float32))
	log.WithField(
		"value", *p.exposure,
	).WithField(
		"min", *p.minExposure,
	).WithField(
		"max", *p.maxExposure,
	).Debug("image exposure changed")
	return nil
}

// saturationChanged is invoked by the device when the picture saturation is
// changed.
func (p *pictureSettingsState) saturationChanged(args []interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.saturation = ptr.ToFloat32(args[0].(float32))
	p.minSaturation = ptr.ToFloat32(args[1].(float32))
	p.maxSaturation = ptr.ToFloat32(args[2].(float32))
	log.WithField(
		"value", *p.saturation,
	).WithField(
		"min", *p.minSaturation,
	).WithField(
		"max", *p.maxSaturation,
	).Debug("image saturation changed")
	return nil
}

// TODO: Implement this
// Title: Timelapse mode
// Description: Timelapse mode.
// Support: 0901;090c;090e
// Triggered: by [SetTimelapseMode](#1-19-4).
// Result:
func (p *pictureSettingsState) timelapseChanged(args []interface{}) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.timeLapseEnabled = ptr.ToBool(args[0].(uint8) == 1)
	p.timeLapseInterval = ptr.ToFloat32(args[1].(float32))
	p.minTimeLapseInterval = ptr.ToFloat32(args[2].(float32))
	p.maxTimeLapseInterval = ptr.ToFloat32(args[3].(float32))
	log.WithField(
		"enabled", args[0].(uint8),
	).WithField(
		"interval", *p.timeLapseInterval,
	).WithField(
		"minInterval", *p.minTimeLapseInterval,
	).WithField(
		"maxInterval", *p.maxTimeLapseInterval,
	).Debug("time lapse changed")
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

func (p *pictureSettingsState) Saturation() (float32, bool) {
	if p.saturation == nil {
		return 0, false
	}
	return *p.saturation, true
}

func (p *pictureSettingsState) MinSaturation() (float32, bool) {
	if p.minSaturation == nil {
		return 0, false
	}
	return *p.minSaturation, true
}

func (p *pictureSettingsState) MaxSaturation() (float32, bool) {
	if p.maxSaturation == nil {
		return 0, false
	}
	return *p.maxSaturation, true
}

func (p *pictureSettingsState) Exposure() (float32, bool) {
	if p.exposure == nil {
		return 0, false
	}
	return *p.exposure, true
}

func (p *pictureSettingsState) MinExposure() (float32, bool) {
	if p.minExposure == nil {
		return 0, false
	}
	return *p.minExposure, true
}

func (p *pictureSettingsState) MaxExposure() (float32, bool) {
	if p.maxExposure == nil {
		return 0, false
	}
	return *p.maxExposure, true
}

func (p *pictureSettingsState) TimeLapseEnabled() (bool, bool) {
	if p.timeLapseEnabled == nil {
		return false, false
	}
	return *p.timeLapseEnabled, true
}

func (p *pictureSettingsState) TimeLapseInterval() (float32, bool) {
	if p.timeLapseInterval == nil {
		return 0, false
	}
	return *p.timeLapseInterval, true
}

func (p *pictureSettingsState) MinTimeLapseInterval() (float32, bool) {
	if p.minTimeLapseInterval == nil {
		return 0, false
	}
	return *p.minTimeLapseInterval, true
}

func (p *pictureSettingsState) MaxTimeLapseInterval() (float32, bool) {
	if p.maxTimeLapseInterval == nil {
		return 0, false
	}
	return *p.maxTimeLapseInterval, true
}
