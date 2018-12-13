package common

import (
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Feature ...
// TODO: Document this
type Feature interface {
	arcommands.D2CFeature
	AccessoryState() AccessoryState
	AnimationsState() AnimationsState
	ARLibsVersionsState() ARLibsVersionsState
	AudioState() AudioState
	CalibrationState() CalibrationState
	CameraSettingsState() CameraSettingsState
	ChargerState() ChargerState
	CommonState() CommonState
	FlightPlanEvent() FlightPlanEvent
	FlightPlanSettingsState() FlightPlanSettingsState
	FlightPlanState() FlightPlanState
	HeadlightsState() HeadlightsState
	MavlinkState() MavlinkState
	NetworkEvent() NetworkEvent
	OverHeatState() OverHeatState
	RunState() RunState
	SettingsState() SettingsState
	WifiSettingsState() WifiSettingsState
}

type feature struct {
	accessoryState          *accessoryState
	animationsState         *animationsState
	arLibsVersionsState     *arLibsVersionsState
	audioState              *audioState
	calibrationState        *calibrationState
	cameraSettingsState     *cameraSettingsState
	chargerState            *chargerState
	commonState             *commonState
	flightPlanEvent         *flightPlanEvent
	flightPlanSettingsState *flightPlanSettingsState
	flightPlanState         *flightPlanState
	headlightsState         *headlightsState
	mavlinkState            *mavlinkState
	networkEvent            *networkEvent
	overHeatState           *overHeatState
	runState                *runState
	settingsState           *settingsState
	wifiSettingsState       *wifiSettingsState
}

// NewFeature ...
// TODO: Document this
func NewFeature() Feature {
	return &feature{
		accessoryState:          &accessoryState{},
		animationsState:         &animationsState{},
		arLibsVersionsState:     &arLibsVersionsState{},
		audioState:              &audioState{},
		calibrationState:        &calibrationState{},
		cameraSettingsState:     &cameraSettingsState{},
		chargerState:            &chargerState{},
		commonState:             &commonState{},
		flightPlanEvent:         &flightPlanEvent{},
		flightPlanSettingsState: &flightPlanSettingsState{},
		flightPlanState:         &flightPlanState{},
		headlightsState:         &headlightsState{},
		mavlinkState:            &mavlinkState{},
		networkEvent:            &networkEvent{},
		overHeatState:           &overHeatState{},
		runState:                &runState{},
		settingsState:           &settingsState{},
		wifiSettingsState:       &wifiSettingsState{},
	}
}

func (f *feature) ID() uint8 {
	return 0
}

func (f *feature) Name() string {
	return "common"
}

// TODO: Add stuff!
func (f *feature) D2CClasses() []arcommands.D2CClass {
	return []arcommands.D2CClass{
		f.accessoryState,
		f.animationsState,
		f.arLibsVersionsState,
		f.audioState,
		f.calibrationState,
		f.cameraSettingsState,
		f.chargerState,
		f.commonState,
		f.flightPlanEvent,
		f.flightPlanSettingsState,
		f.flightPlanState,
		f.headlightsState,
		f.mavlinkState,
		f.networkEvent,
		f.overHeatState,
		f.runState,
		f.settingsState,
		f.wifiSettingsState,
	}
}

func (f *feature) AccessoryState() AccessoryState {
	return f.accessoryState
}

func (f *feature) AnimationsState() AnimationsState {
	return f.animationsState
}

func (f *feature) ARLibsVersionsState() ARLibsVersionsState {
	return f.arLibsVersionsState
}

func (f *feature) AudioState() AudioState {
	return f.audioState
}

func (f *feature) CalibrationState() CalibrationState {
	return f.calibrationState
}

func (f *feature) CameraSettingsState() CameraSettingsState {
	return f.cameraSettingsState
}

func (f *feature) ChargerState() ChargerState {
	return f.chargerState
}

func (f *feature) CommonState() CommonState {
	return f.commonState
}

func (f *feature) FlightPlanEvent() FlightPlanEvent {
	return f.flightPlanEvent
}

func (f *feature) FlightPlanSettingsState() FlightPlanSettingsState {
	return f.flightPlanSettingsState
}

func (f *feature) FlightPlanState() FlightPlanState {
	return f.flightPlanState
}

func (f *feature) HeadlightsState() HeadlightsState {
	return f.headlightsState
}

func (f *feature) MavlinkState() MavlinkState {
	return f.mavlinkState
}

func (f *feature) NetworkEvent() NetworkEvent {
	return f.networkEvent
}

func (f *feature) OverHeatState() OverHeatState {
	return f.overHeatState
}

func (f *feature) RunState() RunState {
	return f.RunState
}

func (f *feature) SettingsState() SettingsState {
	return f.settingsState
}

func (f *feature) WifiSettingsState() WifiSettingsState {
	return f.wifiSettingsState
}
