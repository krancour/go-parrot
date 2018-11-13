package ardrone3

import (
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Feature ...
// TODO: Document this
type Feature interface {
	arcommands.D2CFeature
	PilotingState() PilotingState
	CameraState() CameraState
	GPSState() GPSState
}

type feature struct {
	pilotingState *pilotingState
	cameraState   *cameraState
	gpsState      *gpsState
}

// NewFeature ...
// TODO: Document this
func NewFeature() Feature {
	return &feature{
		pilotingState: &pilotingState{},
		cameraState:   &cameraState{},
		gpsState:      &gpsState{},
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
		f.pilotingState,
		f.cameraState,
		f.gpsState,
	}
}

func (f *feature) PilotingState() PilotingState {
	return f.pilotingState
}

func (f *feature) CameraState() CameraState {
	return f.cameraState
}

func (f *feature) GPSState() GPSState {
	return f.gpsState
}
