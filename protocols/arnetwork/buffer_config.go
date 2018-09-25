package arnetwork

import (
	"fmt"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

// BaseBufferConfig represents configuration common to all buffers, regardless
// of whether they contain inbound or outbound frames.
// nolint: lll
type BaseBufferConfig struct {
	ID            uint8                 // Buffer ID 0 - 255
	FrameType     arnetworkal.FrameType // Type of data
	Size          int32                 // Size of the internal fifo
	MaxDataSize   int32                 // Maximum size of an element in the fifo
	IsOverwriting bool                  // What to do when data is received and the fifo is full
}

// validate validates buffer configuration. This is used internally to
// assert the reasonability of a configuration before attempting to use
// it to initialize a new buffer.
func (b BaseBufferConfig) validate() error {
	if b.FrameType != arnetworkal.FrameTypeAck &&
		b.FrameType != arnetworkal.FrameTypeData &&
		b.FrameType != arnetworkal.FrameTypeLowLatencyData &&
		b.FrameType != arnetworkal.FrameTypeDataWithAck {
		return fmt.Errorf(
			"buffer %d defined with invalid frame type %d",
			b.ID,
			b.FrameType,
		)
	}
	return nil
}
