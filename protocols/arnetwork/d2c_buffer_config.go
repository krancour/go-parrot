package arnetwork

import (
	"fmt"

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
			"input buffer %d defined with invalid frame type %d",
			d.ID,
			d.FrameType,
		)
	}
	return nil
}
