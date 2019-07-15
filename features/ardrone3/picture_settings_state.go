package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Photo settings state from product

const (
	PictureFormatRaw         int32 = 0
	PictureFormatJPEG        int32 = 1
	PictureFormatSnapshot    int32 = 2
	PictureFormatJPEGFisheye int32 = 3

	WhiteBalanceModeAuto      int32 = 0
	WhiteBalanceModeTungsten  int32 = 1
	WhiteBalanceModeDaylight  int32 = 2
	WhiteBalanceModeCloudy    int32 = 3
	WhiteBalanceModeCoolWhite int32 = 4

	VideoStabilizationModeRollPitch int32 = 0
	VideoStabilizationModePitch     int32 = 1
	VideoStabilizationModeRoll      int32 = 2
	VideoStabilizationModeRollNone  int32 = 3

	VideoRecordingModeQuality int32 = 0
	VideoRecordingModeTime    int32 = 1

	VideoFramerate24FPS int32 = 0 // 23.976 frames per second
	VideoFramerate25FPS int32 = 1 // 25 frames per second
	VideoFramerate30FPS int32 = 2 // 29.97 frames per second

	VideoResolutionRec1080Stream480 int32 = 0
	VideoResolutionRec720Stream720  int32 = 1
)

// PictureSettingsState ...
// TODO: Document this
type PictureSettingsState interface {
	lock.ReadLockable
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
	// WhiteBalanceMode returns the white balance mode. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	WhiteBalanceMode() (int32, bool)
	// VideoAutorecordingEnabled returns a boolean value indicating whether video
	// is automatically being captured. A second boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	VideoAutorecordingEnabled() (bool, bool)
	// VideoAutorecordingMassStorageID returns the mass storage ID for any video
	// that is automatically being captured. A boolean value is also returned,
	// indicating whether the first value was reported by the device (true) or a
	// default value (false). This permits callers to distinguish real zero values
	// from default zero values.
	VideoAutorecordingMassStorageID() (uint8, bool)
	// VideoStabilizationMode returns the video stabilization mode. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	VideoStabilizationMode() (int32, bool)
	// VideoRecordingMode returns the video recording mode. A boolean value is
	// also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	VideoRecordingMode() (int32, bool)
	// VideoFramerate returns the video framerate. Be careful! This value maps to
	// a constant and does not directly represent the number of frames per second.
	// A boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	VideoFramerate() (int32, bool)
	// VideoResolutions returns the recording and streaming resolution mode. Be
	// careful! This value maps to a constant and does not directly represent
	// either of these resolutions. A boolean value is also returned, indicating
	// whether the first value was reported by the device (true) or a default
	// value (false). This permits callers to distinguish real zero values from
	// default zero values.
	VideoResolutions() (int32, bool)
}

type pictureSettingsState struct {
	sync.RWMutex
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
	// whiteBalanceMode is the white balance mode
	whiteBalanceMode *int32
	// videoAutorecordingEnabled indicated whether video is automatically captured
	videoAutorecordingEnabled *bool
	// videoAutorecordingMassStorageID is the mass storage id for any video that
	// is automatically captured
	videoAutorecordingMassStorageID *uint8
	// videoStabilizationMode is the video stabilization mode
	videoStabilizationMode *int32
	// videoRecordingMode is the video recording mode
	videoRecordingMode *int32
	// videoFramerate is the video framerate
	videoFramerate *int32
	// videoResolutions represents both the recording an streaming resolution
	videoResolutions *int32
}

func newPictureSettingsState() *pictureSettingsState {
	return &pictureSettingsState{}
}

func (p *pictureSettingsState) ID() uint8 {
	return 20
}

func (p *pictureSettingsState) Name() string {
	return "PictureSettingsState"
}

func (p *pictureSettingsState) D2CCommands(
	log *log.Entry,
) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"PictureFormatChanged",
			[]interface{}{
				int32(0), // type,
			},
			p.pictureFormatChanged,
			log,
		),
		arcommands.NewD2CCommand(
			1,
			"AutoWhiteBalanceChanged",
			[]interface{}{
				int32(0), // type,
			},
			p.autoWhiteBalanceChanged,
			log,
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
			log,
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
			log,
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
			log,
		),
		arcommands.NewD2CCommand(
			5,
			"VideoAutorecordChanged",
			[]interface{}{
				uint8(0), // enabled,
				uint8(0), // mass_storage_id,
			},
			p.videoAutorecordChanged,
			log,
		),
		arcommands.NewD2CCommand(
			6,
			"VideoStabilizationModeChanged",
			[]interface{}{
				int32(0), // mode,
			},
			p.videoStabilizationModeChanged,
			log,
		),
		arcommands.NewD2CCommand(
			7,
			"VideoRecordingModeChanged",
			[]interface{}{
				int32(0), // mode,
			},
			p.videoRecordingModeChanged,
			log,
		),
		arcommands.NewD2CCommand(
			8,
			"VideoFramerateChanged",
			[]interface{}{
				int32(0), // framerate,
			},
			p.videoFramerateChanged,
			log,
		),
		arcommands.NewD2CCommand(
			9,
			"VideoResolutionsChanged",
			[]interface{}{
				int32(0), // type,
			},
			p.videoResolutionsChanged,
			log,
		),
	}
}

// pictureFormatChanged is invoked by the device when the picture format is
// changed.
func (p *pictureSettingsState) pictureFormatChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.format = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"type", *p.format,
	).Debug("picture format changed")
	return nil
}

// autoWhiteBalanceChanged is invoked by the device when the white balance mode
// is changed.
func (p *pictureSettingsState) autoWhiteBalanceChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.whiteBalanceMode = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"type", *p.whiteBalanceMode,
	).Debug("white balance changed")
	return nil
}

// expositionChanged is invoked by the device when exposure is changed.
func (p *pictureSettingsState) expositionChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
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
func (p *pictureSettingsState) saturationChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
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

// timelapseChanged is invoked by the device when time lapse photography
// settings are changed.
func (p *pictureSettingsState) timelapseChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
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

// videoAutorecordChanged is invoked by the device when video autorecording is
// enabled or disabled.
func (p *pictureSettingsState) videoAutorecordChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.videoAutorecordingEnabled = ptr.ToBool(args[0].(uint8) == 1)
	p.videoAutorecordingMassStorageID = ptr.ToUint8(args[1].(uint8))
	log.WithField(
		"enabled", args[0].(uint8),
	).WithField(
		"mass_storage_id", *p.videoAutorecordingMassStorageID,
	).Debug("video autorecording changed")
	return nil
}

// videoStabilizationModeChanged is invoked by the device when the video
// stabilization mode is changed.
func (p *pictureSettingsState) videoStabilizationModeChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.videoStabilizationMode = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"mode", *p.videoStabilizationMode,
	).Debug("video stabilization mode changed")
	return nil
}

// videoRecordingModeChanged is invoked by the device when the video recording
// mode is changed.
func (p *pictureSettingsState) videoRecordingModeChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.videoRecordingMode = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"mode", *p.videoRecordingMode,
	).Debug("video recording mode changed")
	return nil
}

// videoFramerateChanged is invoked by the devide when the video framerate is
// changed.
func (p *pictureSettingsState) videoFramerateChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.videoFramerate = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"framerate", *p.videoFramerate,
	).Debug("video framerate changed")
	return nil
}

// videoResolutionsChanged is invoked by the device when the video resolution
// is changed.
func (p *pictureSettingsState) videoResolutionsChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.videoResolutions = ptr.ToInt32(args[0].(int32))
	log.WithField(
		"type", *p.videoResolutions,
	).Debug("video resolutions changed")
	return nil
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

func (p *pictureSettingsState) WhiteBalanceMode() (int32, bool) {
	if p.whiteBalanceMode == nil {
		return 0, false
	}
	return *p.whiteBalanceMode, true
}

func (p *pictureSettingsState) VideoAutorecordingEnabled() (bool, bool) {
	if p.videoAutorecordingEnabled == nil {
		return false, false
	}
	return *p.videoAutorecordingEnabled, true
}

func (p *pictureSettingsState) VideoAutorecordingMassStorageID() (uint8, bool) {
	if p.videoAutorecordingMassStorageID == nil {
		return 0, false
	}
	return *p.videoAutorecordingMassStorageID, true
}

func (p *pictureSettingsState) VideoStabilizationMode() (int32, bool) {
	if p.videoStabilizationMode == nil {
		return 0, false
	}
	return *p.videoStabilizationMode, true
}

func (p *pictureSettingsState) VideoRecordingMode() (int32, bool) {
	if p.videoRecordingMode == nil {
		return 0, false
	}
	return *p.videoRecordingMode, true
}

func (p *pictureSettingsState) VideoFramerate() (int32, bool) {
	if p.videoFramerate == nil {
		return 0, false
	}
	return *p.videoFramerate, true
}

func (p *pictureSettingsState) VideoResolutions() (int32, bool) {
	if p.videoResolutions == nil {
		return 0, false
	}
	return *p.videoResolutions, true
}
