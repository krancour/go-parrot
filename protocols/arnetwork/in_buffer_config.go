package arnetwork

import (
	"fmt"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

// InBufferConfig represents the configuration of a buffer for inbound frames.
// nolint: lll
type InBufferConfig struct {
	ID            uint8                 // Buffer ID 0 - 255
	FrameType     arnetworkal.FrameType // Type of data
	Size          int32                 // Size of the internal fifo
	MaxDataSize   int32                 // Maximum size of an element in the fifo
	IsOverwriting bool                  // What to do when data is received and the fifo is full
}

// validate validates buffer configuration. This is used internally to
// assert the reasonability of a configuration before attempting to use
// it to initialize a new buffer.
func (i InBufferConfig) validate() error {
	if i.FrameType != arnetworkal.FrameTypeAck &&
		i.FrameType != arnetworkal.FrameTypeData &&
		i.FrameType != arnetworkal.FrameTypeLowLatencyData &&
		i.FrameType != arnetworkal.FrameTypeDataWithAck {
		return fmt.Errorf(
			"input buffer %d defined with invalid frame type %d",
			i.ID,
			i.FrameType,
		)
	}
	return nil
}
