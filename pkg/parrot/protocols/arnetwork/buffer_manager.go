package arnetwork

import (
	"github.com/krancour/drone-examples/pkg/parrot/protocols/arnetworkal"
)

// BufferManager ...
// TODO: Document this
type BufferManager interface {
	Send(Frame) error
	// TODO: Not sure what a receive from this thing should look like.
}

type bufferManager struct {
	outBuffers map[uint8]*Buffer
	inBuffers  map[uint8]*Buffer
}

// NewBufferManager ...
// TODO: Document this
func NewBufferManager(
	conn arnetworkal.Connection,
	buffers ...*Buffer,
) (BufferManager, error) {
	// TODO: Add buffers to bufferManager
	b := &bufferManager{}
	for _, buffer := range buffers {
		if err := buffer.validate(); err != nil {
			return nil, err
		}
		// Automatically create an ack buffer for every buffer...
		// nolint: lll
		ackBuffer := &Buffer{
			ID:              buffer.ID + ackBufferOffset,
			DataType:        arnetworkal.FrameTypeAck,
			SendingWaitTime: 0, // Deprecated
			AckTimeout:      0, // Unused
			NumberOfRetries: 0, // Unused
			NumberOfCells:   1, // Never more than one ack at a time
			DataCopyMaxSize: 1, // One byte of data: the sequence number
			IsOverwriting:   0, // Useless by design: there is only one ack waiting at a time
		}
		if buffer.Dir == BufferDirectionOutbound {
			// Hang on to this buffer...
			b.outBuffers[buffer.ID] = buffer
			// And the corresponding ack buffer
			ackBuffer.Dir = BufferDirectionInbound
			b.inBuffers[ackBuffer.ID] = ackBuffer
		} else {
			// Hang on to this buffer...
			b.inBuffers[buffer.ID] = buffer
			ackBuffer.Dir = BufferDirectionOutbound
			// And the corresponding ack buffer
			b.outBuffers[ackBuffer.ID] = ackBuffer
		}
	}
	return b, nil
}

// Send ...
// TODO: Document this
func (b *bufferManager) Send(frame Frame) error {
	// TODO: Implement this
	// The logic required here is a little bit complex
	return nil
}
