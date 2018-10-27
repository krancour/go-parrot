package arnetwork

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

// D2CBufferConfig represents the configuration of a buffer for frame sent from
// the device to the client.
// nolint: lll
type D2CBufferConfig struct {
	ID            uint8                 // Buffer ID 0 - 255
	FrameType     arnetworkal.FrameType // Type of data
	Size          int32                 // Size of the internal fifo
	MaxDataSize   int32                 // Maximum size of an element in the fifo
	IsOverwriting bool                  // What to do when data is received and the fifo is full
}

// validate validates buffer configuration. This is used internally to
// assert the reasonability of a configuration before attempting to use
// it to initialize a new buffer.
func (d D2CBufferConfig) validate() error {
	if d.FrameType != arnetworkal.FrameTypeAck &&
		d.FrameType != arnetworkal.FrameTypeData &&
		d.FrameType != arnetworkal.FrameTypeLowLatencyData &&
		d.FrameType != arnetworkal.FrameTypeDataWithAck {
		return fmt.Errorf(
			"d2c buffer %d defined with invalid frame type %d",
			d.ID,
			d.FrameType,
		)
	}
	if d.Size < 1 {
		return fmt.Errorf(
			"d2c buffer %d defined with invalid size %d",
			d.ID,
			d.Size,
		)
	}
	log.WithField("id", d.ID).Debug("d2c buffer config is valid")
	return nil
}
