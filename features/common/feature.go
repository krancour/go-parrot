package common

import (
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Feature ...
// TODO: Document this
type Feature interface {
	arcommands.D2CFeature
	CommonState() CommonState
}

type feature struct {
	commonState *commonState
}

// NewFeature ...
// TODO: Document this
func NewFeature() Feature {
	return &feature{
		commonState: &commonState{},
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
		f.commonState,
	}
}

func (f *feature) CommonState() CommonState {
	return f.commonState
}
