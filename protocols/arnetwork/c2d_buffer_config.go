package arnetwork

import (
	"fmt"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

// C2DBufferConfig represents the configuration of a buffer for frames being
// sent from the client to the device.
// nolint: lll
type C2DBufferConfig struct {
	ID            uint8                 // Buffer ID 0 - 255
	FrameType     arnetworkal.FrameType // Type of data
	Size          int32                 // Size of the internal fifo
	MaxDataSize   int32                 // Maximum size of an element in the fifo
	IsOverwriting bool                  // What to do when data is received and the fifo is full
	AckTimeout    time.Duration         // Time before considering a frame lost
	MaxRetries    int                   // Number of retries before considering a frame lost
}

// validate validates buffer configuration. This is used internally to
// assert the reasonability of a configuration before attempting to use
// it to initialize a new buffer.
func (c C2DBufferConfig) validate() error {
	if c.FrameType != arnetworkal.FrameTypeAck &&
		c.FrameType != arnetworkal.FrameTypeData &&
		c.FrameType != arnetworkal.FrameTypeLowLatencyData &&
		c.FrameType != arnetworkal.FrameTypeDataWithAck {
		return fmt.Errorf(
			"c2d buffer %d defined with invalid frame type %d",
			c.ID,
			c.FrameType,
		)
	}
	return nil
}