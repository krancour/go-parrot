package animation

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

type Feature interface {
	arcommands.D2CFeature
}

type feature struct {
}

func NewFeature(c2dCommandClient arcommands.C2DCommandClient) Feature {
	return &feature{}
}

func (f *feature) FeatureID() uint8 {
	return 144
}

func (f *feature) FeatureName() string {
	return "animation"
}

func (f *feature) D2CClasses() []arcommands.D2CClass {
	return []arcommands.D2CClass{f}
}

func (f *feature) ClassID() uint8 {
	return 0
}

func (f *feature) ClassName() string {
	return ""
}

func (f *feature) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{}
}
