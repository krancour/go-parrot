package arnetwork

import (
	"container/ring"
	"fmt"
	"sync"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
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

// BufferConfig represents the configuration for a buffer. This is exported,
// so can be instantiated at higher levels of abstraction and passed to the
// NewBufferManager(...) function, which will use the configuration in
// initializing buffers implemented by an unexported type.
// nolint: lll
type BufferConfig struct {
	Direction     BufferDirection       // Inbound or outbound
	ID            uint8                 // Buffer ID 0 - 255
	DataType      arnetworkal.FrameType // Type of data
	AckTimeout    time.Duration         // Time before considering a frame lost
	MaxRetries    int                   // Number of retries before considering a frame lost
	Size          int32                 // Size of the internal fifo
	MaxDataSize   int32                 // Maximum size of an element in the fifo
	IsOverwriting bool                  // What to do when data is received and the fifo is full
}

// type buffer interface {
// 	write(Frame) bool
// 	read() (Frame, bool, error)
// }

// buffer implements a ring buffer. Behavior when writing to a full buffer is
// determined by BufferConfig.IsOverwriting. Outbound buffers continuously
// read from their own head to send frames out via the arnetworkal connection,
// while inbound buffers continuously read from their own head to invoke the
// specified function.
type buffer struct {
	BufferConfig
	mtx  sync.Mutex
	seq  uint8 // Internal sequence number tracking
	head *ring.Ring
	next *ring.Ring
	conn arnetworkal.Connection
}

// validate validates buffer configuration. This is used internally to
// assert the reasonability of a configuration before attempting to use
// it to initialize a new buffer.
func (b BufferConfig) validate() error {
	if b.Direction != BufferDirectionOutbound &&
		b.Direction != BufferDirectionInbound {
		return fmt.Errorf(
			"buffer %d defined with invalid direction %d",
			b.ID,
			b.Direction,
		)
	}
	if b.DataType != arnetworkal.FrameTypeAck &&
		b.DataType != arnetworkal.FrameTypeData &&
		b.DataType != arnetworkal.FrameTypeLowLatencyData &&
		b.DataType != arnetworkal.FrameTypeDataWithAck {
		return fmt.Errorf(
			"buffer %d defined with invalid direction %d",
			b.ID,
			b.Direction,
		)
	}
	return nil
}

func newBuffer(bufCfg BufferConfig, conn arnetworkal.Connection) *buffer {
	buf := &buffer{
		BufferConfig: bufCfg,
		next:         ring.New(int(bufCfg.Size)),
		conn:         conn,
	}
	// TODO: For outbound buffers, start a goroutine that reads from the
	// head of the buffer and writes out arnetworkal frames.
	// TODO: For inbound buffers, start a goroutine that reads from the head
	// of the buffer and does ???
	return buf
}

func (b *buffer) write(frame Frame) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if b.next == b.head {
		// The next place to write to is the head. i.e. The buffer is full.
		// How do we deal with it?
		// Are we "overwriting?"
		if !b.IsOverwriting {
			// We're not "overwriting," so we'll just drop the new frame by
			// simply returning.
			return
		}
		// We are "overwriting," so we'll advance the head (i.e. drop oldest
		// frame in the buffer.
		b.head = b.head.Next()
	}
	b.next.Value = frame
	// TODO: We may need to send an ack here, depending on direction and type!
	if b.head == nil {
		b.head = b.next
	}
	b.next = b.next.Next()
}

// read reads a single frame from the head of the buffer. A boolean is also
// returned to indicate success or failure. It is possible, for instance, for a
// read ot return nothing without the occurence of an error, simply because the
// buffer is empty. An error is also potentially returned.
func (b *buffer) read() (Frame, bool, error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if b.head == nil {
		// The buffer is empty.
		return Frame{}, false, nil
	}
	frame, ok := b.head.Value.(Frame)
	b.head = b.head.Next()
	// If we've advanced the head to the next spot we will write to, we have an
	// empty buffer and we should represent that by pointing the head to nil.
	if b.head == b.next {
		b.head = nil
	}
	if !ok {
		return frame,
			false,
			fmt.Errorf("data read from buffer was not an arnetwork frame")
	}
	return frame, true, nil
}

// TODO: I might want to make two distinc types of buffer-- InBuffer and
// OutBuffer, with both composed of a "base" Buffer + attributes specific to
// the direction. For instance, for an InBuffer, I might wish to pass in a
// func as a callback.
