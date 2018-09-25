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
	outBuffers map[uint8]*outBuffer
	inBuffers  map[uint8]*inBuffer
}

// NewBufferManager returns a new BufferManager.
func NewBufferManager(
	conn arnetworkal.Connection,
	inBufCfgs []InBufferConfig,
	outBufCfgs []OutBufferConfig,
) (BufferManager, error) {
	b := &bufferManager{
		conn: conn,
	}

	for _, bufCfg := range inBufCfgs {
		if err := bufCfg.validate(); err != nil {
			return nil, err
		}
		b.inBuffers[bufCfg.ID] = newInBuffer(bufCfg)
		if bufCfg.FrameType == arnetworkal.FrameTypeDataWithAck {
			// Automatically create an ack buffer...
			// nolint: lll
			b.outBuffers[bufCfg.ID] = newOutBuffer(
				OutBufferConfig{
					BaseBufferConfig: BaseBufferConfig{
						ID:            bufCfg.ID + ackBufferOffset,
						FrameType:     arnetworkal.FrameTypeAck,
						Size:          1,     // Never more than one ack at a time
						MaxDataSize:   1,     // One byte of data: the sequence number
						IsOverwriting: false, // Useless by design: there is only one ack waiting at a time
					},
					AckTimeout: 0, // Unused
					MaxRetries: 0, // Unused
				},
				conn,
			)
		}
	}

	for _, bufCfg := range outBufCfgs {
		if err := bufCfg.validate(); err != nil {
			return nil, err
		}
		b.outBuffers[bufCfg.ID] = newOutBuffer(bufCfg, conn)
		// Automatically create an ack buffer...
		// nolint: lll
		b.inBuffers[bufCfg.ID] = newInBuffer(
			InBufferConfig{
				BaseBufferConfig: BaseBufferConfig{
					ID:            bufCfg.ID + ackBufferOffset,
					FrameType:     arnetworkal.FrameTypeAck,
					Size:          1,     // Never more than one ack at a time
					MaxDataSize:   1,     // One byte of data: the sequence number
					IsOverwriting: false, // Useless by design: there is only one ack waiting at a time
				},
				// TODO: Not sure yet how to connect this to the function
				// that writes frames out and is expecting an ack.
				CallBack: func(frame Frame) {

				},
			},
		)
	}

	go b.receiveFrames()

	return b, nil
}

// Send adds the provided frame to the correct buffer and returns an error if
// the correct buffer does not exist. This function does not immediately result
// in the frame being sent over the network. A separate goroutine will run in
// the background to drain each buffer.
func (b *bufferManager) Send(frame Frame) error {
	buffer, ok := b.outBuffers[frame.BufferID]
	if !ok {
		return fmt.Errorf(
			"error sending frame: unknown buffer %d",
			frame.BufferID,
		)
	}
	buffer.write(frame)
	return nil
}

func (b *bufferManager) receiveFrames() {
	for {
		netFrame, ok, err := b.conn.Receive()
		if err != nil {
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
			BufferID: netFrame.ID,
			seq:      netFrame.Seq,
			Data:     netFrame.Data,
		}
		// Put the frame on the correct buffer...
		buf, ok := b.inBuffers[frame.BufferID]
		if !ok {
			// No buffer found to receive this frame.
			log.Printf(
				"received frame for unknown buffer %d\n",
				frame.BufferID,
			)
			// krancour: Pretty sure we're not supposed to ack frames that we
			// aren't even going to accept into some buffer. We can revisit that
			// later if this is called into question.
			continue
		}
		ok = buf.write(frame)
		// If the frame was accepted onto the buffer and acknowledgement was
		// requested, send it...
		if ok && netFrame.Type == arnetworkal.FrameTypeDataWithAck {
			ackBufID := netFrame.ID + ackBufferOffset
			ackBuf, ok := b.outBuffers[ackBufID]
			if !ok {
				log.Printf(
					"error sending ack: outbound buffer %d does not exist\n",
					ackBufID,
				)
			}
			ackBuf.write(Frame{
				BufferID: ackBufID,
				Data:     []byte(fmt.Sprintf("%d", netFrame.Seq)),
			})
		}
	}
}
