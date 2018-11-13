package arcommands

import (
	"bytes"
	"encoding/binary"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetwork"
)

// D2CCommandServer ...
// TODO: Document this
type D2CCommandServer interface {
	Start()
}

type d2cCommandServer struct {
	d2cChs      map[uint8]<-chan arnetwork.Frame
	d2cFeatures map[uint8]D2CFeature
}

// NewD2CCommandServer ...
// TODO: Document this
func NewD2CCommandServer(
	d2cChs map[uint8]<-chan arnetwork.Frame,
	d2cFeatures []D2CFeature,
) D2CCommandServer {
	d2cFeatMap := map[uint8]D2CFeature{}
	for _, d2cFeature := range d2cFeatures {
		d2cFeatMap[d2cFeature.ID()] = d2cFeature
	}
	return &d2cCommandServer{
		d2cChs:      d2cChs,
		d2cFeatures: d2cFeatMap,
	}
}

// Run ...
// TODO: Document this
// TODO: Return errors? Or use an error channel?
func (d *d2cCommandServer) Start() {
	for bufID, d2cCh := range d.d2cChs {
		go d.receiveCommands(bufID, d2cCh)
	}
}

// TODO: Move this into a separate file
func (d *d2cCommandServer) receiveCommands(
	_ uint8,
	d2cCh <-chan arnetwork.Frame,
) {
	for frame := range d2cCh {
		featureID, classID, commandID, err := parseIDS(frame.Data)
		if err != nil {
			log.Error(err)
			continue
		}
		feature, ok := d.d2cFeatures[featureID]
		if !ok {
			log.WithField(
				"featureID", featureID,
			).Warn("feature not found")
			continue
		}
		class, ok := feature.D2CClasses()[classID]
		if !ok {
			log.WithField(
				"featureID", featureID,
			).WithField(
				"classID", classID,
			).Warn("class not found")
			continue
		}
		command, ok := class.d2cCommands()[commandID]
		if !ok {
			log.WithField(
				"featureID", featureID,
			).WithField(
				"classID", classID,
			).WithField(
				"commandID", commandID,
			).Warn("command not found")
			continue
		}
		if err := command.execute(frame.Data); err != nil {
			log.Error(err)
		}
	}
}

func parseIDS(data []byte) (
	featureID uint8,
	classID uint8,
	commandID uint16,
	err error,
) {
	buf := bytes.NewBuffer(data)
	if err = binary.Read(buf, binary.LittleEndian, &featureID); err != nil {
		return
	}
	if err = binary.Read(buf, binary.LittleEndian, &classID); err != nil {
		return
	}
	err = binary.Read(buf, binary.LittleEndian, &commandID)
	return
}
