package arnetwork

import (
	"fmt"
	"time"

	"github.com/krancour/drone-examples/pkg/parrot/protocols/arnetworkal"
)

// BufferDirection is a type for constants used to indicate the direction of
// a buffer-- outbound or inbound.
type BufferDirection uint8

const (
	// bufferDirectionUnknown is the zero value for type BufferDirection and
	// represents an unspecified buffer direction. This should not be
	// deliberately used!
	bufferDirectionUnknown BufferDirection = 0
	// BufferDirectionOutbound represents a buffer used for sending data.
	BufferDirectionOutbound BufferDirection = 1
	// BufferDirectionInbound represents a buffer used for receiving data.
	BufferDirectionInbound BufferDirection = 2

	ackBufferOffset uint8 = 128
)

// Buffer ...
// TODO: Document this
// nolint: lll
type Buffer struct {
	Dir             BufferDirection       // Inbound or outbound
	ID              uint8                 // Buffer ID 0 - 255
	DataType        arnetworkal.FrameType // Type of data
	SendingWaitTime time.Duration         // Deprecated TODO: Probably remove this
	AckTimeout      time.Duration         // Time before considering a frame lost
	NumberOfRetries int                   // Number of retries before considering a frame lost
	NumberOfCells   int32                 // Size of the internal fifo
	DataCopyMaxSize int32                 // Maximum size of an element in the fifo
	IsOverwriting   int                   // What to do when data is received and the fifo is full
	seq             uint8                 // Internal sequence number tracking
}

func (b *Buffer) validate() error {
	if b.Dir != BufferDirectionOutbound &&
		b.Dir != BufferDirectionInbound {
		return fmt.Errorf(
			"buffer %d defined with invalid direction %d",
			b.ID,
			b.Dir,
		)
	}
	if b.DataType != arnetworkal.FrameTypeAck &&
		b.DataType != arnetworkal.FrameTypeData &&
		b.DataType != arnetworkal.FrameTypeLowLatencyData &&
		b.DataType != arnetworkal.FrameTypeDataWithAck {
		return fmt.Errorf(
			"buffer %d defined with invalid direction %d",
			b.ID,
			b.Dir,
		)
	}
	return nil
}

// TODO: I might want to make two distinc types of buffer-- InBuffer and
// OutBuffer, with both composed of a "base" Buffer + attributes specific to
// the direction. For instance, for an InBuffer, I might wish to pass in a
// func as a callback.
