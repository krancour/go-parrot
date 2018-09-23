package arnetwork

import (
	"fmt"
	"log"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

// BufferManager sends and receives arnetwork frames via a series of outbound
// and inbound buffers.
type BufferManager interface {
	Send(Frame) error
	// TODO: Not sure what a receive from this thing should look like.
	// Maybe just allow callbacks of some sort to be registered?
}

type bufferManager struct {
	conn       arnetworkal.Connection
	outBuffers map[uint8]*buffer
	inBuffers  map[uint8]*buffer
}

// NewBufferManager returns a new BufferManager.
func NewBufferManager(
	conn arnetworkal.Connection,
	bufCfgs ...BufferConfig,
) (BufferManager, error) {
	b := &bufferManager{
		conn: conn,
	}
	for _, bufCfg := range bufCfgs {
		if err := bufCfg.validate(); err != nil {
			return nil, err
		}
		buf := newBuffer(bufCfg, conn)
		// Automatically create an ack buffer for every buffer...
		// nolint: lll
		ackBuf := newBuffer(
			BufferConfig{
				ID:            buf.ID + ackBufferOffset,
				DataType:      arnetworkal.FrameTypeAck,
				AckTimeout:    0,     // Unused
				MaxRetries:    0,     // Unused
				Size:          1,     // Never more than one ack at a time
				MaxDataSize:   1,     // One byte of data: the sequence number
				IsOverwriting: false, // Useless by design: there is only one ack waiting at a time
			},
			conn,
		)
		if buf.Direction == BufferDirectionOutbound {
			// Hang on to this buffer...
			b.outBuffers[buf.ID] = buf
			// And the corresponding ack buffer
			ackBuf.Direction = BufferDirectionInbound
			b.inBuffers[ackBuf.ID] = ackBuf
		} else {
			// Hang on to this buffer...
			b.inBuffers[buf.ID] = buf
			ackBuf.Direction = BufferDirectionOutbound
			// And the corresponding ack buffer
			b.outBuffers[ackBuf.ID] = ackBuf
		}
	}

	go b.receiveFrames()

	return b, nil
}

// Send adds the provided frame to the correct buffer and returns an error if
// the correct buffer does not exist. This function does not immediately result
// in the frame being sent over the network. A separate goroutine will run in
// the background to drain each buffer.
func (b *bufferManager) Send(frame Frame) error {
	buffer, ok := b.outBuffers[frame.ID]
	if !ok {
		return fmt.Errorf("error sending frame: unknown buffer %d", frame.ID)
	}
	buffer.write(frame)
	return nil
}

func (b *bufferManager) receiveFrames() {
	for {
		netFrame, ok, err := b.conn.Receive()
		if err != nil {
			// TODO: Can we do anything here besides log this?
			log.Printf("error receiving frame: %s\n", err)
			continue
		}
		if !ok {
			// There was no error, but we also didn't receive any data. This
			// only happens when the underlying arnetworkal connection has been
			// closed, so we'll just return to let this goroutine conclude.
			return
		}
		// If we make it this far, we have a frame!
		// Unpack the arnetworkal frame into an arnetwork frame...
		frame := Frame{
			ID:   netFrame.ID,
			seq:  netFrame.Seq,
			Data: netFrame.Data,
		}
		// Put the frame on the correct buffer...
		buffer, ok := b.inBuffers[frame.ID]
		if !ok {
			// No buffer found to receive this frame.
			log.Printf("received frame for unknown buffer %d\n", frame.ID)
			// krancour: Pretty sure we're not supposed to ack frames that we
			// aren't even going to accept into some buffer. We can revisit that
			// later if this is called into question.
			continue
		}
		buffer.write(frame)
	}
}
