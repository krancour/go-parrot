package arcommands

import (
	"bytes"
	"encoding/binary"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetwork"
	"github.com/pkg/errors"
)

// D2CCommandServer ...
// TODO: Document this
type D2CCommandServer interface {
	Start()
}

type d2cCommandServer struct {
	d2cChs      map[uint8]<-chan arnetwork.Frame
	d2cCommands map[string]D2CCommand
}

// NewD2CCommandServer ...
// TODO: Document this
func NewD2CCommandServer(
	d2cChs map[uint8]<-chan arnetwork.Frame,
	d2cFeatures []D2CFeature,
) (D2CCommandServer, error) {
	d2cCommands := map[string]D2CCommand{}
	for _, feature := range d2cFeatures {
		for _, class := range feature.D2CClasses() {
			for _, command := range class.D2CCommands() {
				key := getCommandKey(feature.ID(), class.ID(), command.ID())
				if _, ok := d2cCommands[key]; ok {
					return nil, errors.Errorf("command with key %s already defined", key)
				}
				d2cCommands[key] = command
			}
		}
	}
	return &d2cCommandServer{
		d2cChs:      d2cChs,
		d2cCommands: d2cCommands,
	}, nil
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
		l := log.WithField(
			"featureID", featureID,
		).WithField(
			"classID", classID,
		).WithField(
			"commandID", commandID,
		)
		key := getCommandKey(featureID, classID, commandID)
		command, ok := d.d2cCommands[key]
		if !ok {
			l.Warn("command not found")
			continue
		}
		if err := command.execute(frame.Data); err != nil {
			l.Error(err)
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
		err = errors.Wrap(err, "error parsing featureID from command")
		return
	}
	if err = binary.Read(buf, binary.LittleEndian, &classID); err != nil {
		err = errors.Wrap(err, "error parsing classID from command")
		return
	}
	if err = binary.Read(buf, binary.LittleEndian, &commandID); err != nil {
		err = errors.Wrap(err, "error parsing commandID from command")
	}
	return
}

func getCommandKey(featureID, classID uint8, commandID uint16) string {
	return fmt.Sprintf("%d:%d:%d", featureID, classID, commandID)
}
