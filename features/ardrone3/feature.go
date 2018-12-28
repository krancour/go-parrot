package ardrone3

import (
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Feature ...
// TODO: Document this
type Feature interface {
	arcommands.D2CFeature
	// AccessoryState() AccessoryState
	AntiflickeringState() AntiflickeringState
	CameraState() CameraState
	GPSSettingsState() GPSSettingsState
	GPSState() GPSState
	MediaRecordEvent() MediaRecordEvent
	MediaRecordState() MediaRecordState
	MediaStreamingState() MediaStreamingState
	NetworkSettingsState() NetworkSettingsState
	NetworkState() NetworkState
	PictureSettingsState() PictureSettingsState
	PilotingEvent() PilotingEvent
	PilotingSettingsState() PilotingSettingsState
	PilotingState() PilotingState
	// PROState() PROState
	SettingsState() SettingsState
	// SoundState() SoundState
	SpeedSettingsState() SpeedSettingsState
}

type feature struct {
	// accessoryState        *accessoryState
	antiflickeringState   *antiflickeringState
	cameraState           *cameraState
	gpsSettingsState      *gpsSettingsState
	gpsState              *gpsState
	mediaRecordEvent      *mediaRecordEvent
	mediaRecordState      *mediaRecordState
	mediaStreamingState   *mediaStreamingState
	networkSettingsState  *networkSettingsState
	networkState          *networkState
	pictureSettingsState  *pictureSettingsState
	pilotingEvent         *pilotingEvent
	pilotingSettingsState *pilotingSettingsState
	pilotingState         *pilotingState
	// proState              *proState
	settingsState *settingsState
	// soundState         *soundState
	speedSettingsState *speedSettingsState
}

// NewFeature ...
// TODO: Document this
func NewFeature() Feature {
	return &feature{
		// accessoryState:        &accessoryState{},
		antiflickeringState:   &antiflickeringState{},
		cameraState:           &cameraState{},
		gpsSettingsState:      &gpsSettingsState{},
		gpsState:              &gpsState{},
		mediaRecordEvent:      &mediaRecordEvent{},
		mediaRecordState:      &mediaRecordState{},
		mediaStreamingState:   &mediaStreamingState{},
		networkSettingsState:  &networkSettingsState{},
		networkState:          &networkState{},
		pictureSettingsState:  &pictureSettingsState{},
		pilotingEvent:         &pilotingEvent{},
		pilotingSettingsState: &pilotingSettingsState{},
		pilotingState:         &pilotingState{},
		// proState:              &proState{},
		settingsState: &settingsState{},
		// soundState:         &soundState{},
		speedSettingsState: &speedSettingsState{},
	}
}

func (f *feature) ID() uint8 {
	return 1
}

func (f *feature) Name() string {
	return "ardrone3"
}

// TODO: Add stuff!
func (f *feature) D2CClasses() []arcommands.D2CClass {
	return []arcommands.D2CClass{
		// f.accessoryState,
		f.antiflickeringState,
		f.cameraState,
		f.gpsSettingsState,
		f.gpsState,
		f.mediaRecordEvent,
		f.mediaRecordState,
		f.mediaStreamingState,
		f.networkSettingsState,
		f.networkState,
		f.pictureSettingsState,
		f.pilotingEvent,
		f.pilotingSettingsState,
		f.pilotingState,
		// f.proState,
		f.settingsState,
		// f.soundState,
		f.speedSettingsState,
	}
}

// func (f *feature) AccessoryState() AccessoryState {
// 	return f.accessoryState
// }

func (f *feature) AntiflickeringState() AntiflickeringState {
	return f.antiflickeringState
}

func (f *feature) CameraState() CameraState {
	return f.cameraState
}

func (f *feature) GPSSettingsState() GPSSettingsState {
	return f.gpsSettingsState
}

func (f *feature) GPSState() GPSState {
	return f.gpsState
}

func (f *feature) MediaRecordEvent() MediaRecordEvent {
	return f.mediaRecordEvent
}

func (f *feature) MediaRecordState() MediaRecordState {
	return f.mediaRecordState
}

func (f *feature) MediaStreamingState() MediaStreamingState {
	return f.mediaStreamingState
}

func (f *feature) NetworkSettingsState() NetworkSettingsState {
	return f.networkSettingsState
}

func (f *feature) NetworkState() NetworkState {
	return f.networkState
}

func (f *feature) PictureSettingsState() PictureSettingsState {
	return f.PictureSettingsState
}

func (f *feature) PilotingEvent() PilotingEvent {
	return f.pilotingEvent
}

func (f *feature) PilotingSettingsState() PilotingSettingsState {
	return f.pilotingSettingsState
}

func (f *feature) PilotingState() PilotingState {
	return f.pilotingState
}

// func (f *feature) PROState() PROState {
// 	return f.proState
// }

func (f *feature) SettingsState() SettingsState {
	return f.settingsState
}

// func (f *feature) SoundState() SoundState {
// 	return f.soundState
// }

func (f *feature) SpeedSettingsState() SpeedSettingsState {
	return f.speedSettingsState
}
