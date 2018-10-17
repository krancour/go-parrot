package arnetwork

import (
	"fmt"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

// OutBufferConfig represents the configuration of a buffer for outbound frames.
// nolint: lll
type OutBufferConfig struct {
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
func (o OutBufferConfig) validate() error {
	if o.FrameType != arnetworkal.FrameTypeAck &&
		o.FrameType != arnetworkal.FrameTypeData &&
		o.FrameType != arnetworkal.FrameTypeLowLatencyData &&
		o.FrameType != arnetworkal.FrameTypeDataWithAck {
		return fmt.Errorf(
			"output buffer %d defined with invalid frame type %d",
			o.ID,
			o.FrameType,
		)
	}
	return nil
}
