package arnetwork

import (
	"log"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

// BufferManager sends and receives arnetwork frames via a series of outbound
// and inbound buffers.
type BufferManager interface {
	OutC(uint8) chan<- Frame
	InC(uint8) <-chan Frame
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
		buf := newInBuffer(bufCfg)
		b.inBuffers[bufCfg.ID] = buf
		if bufCfg.FrameType == arnetworkal.FrameTypeDataWithAck {
			// Automatically create an ack buffer...
			ackBufID := bufCfg.ID + ackBufferOffset
			// nolint: lll
			ackBuf := newOutBuffer(
				OutBufferConfig{
					ID:            ackBufID,
					FrameType:     arnetworkal.FrameTypeAck,
					Size:          1,     // Never more than one ack at a time
					MaxDataSize:   1,     // One byte of data: the sequence number
					IsOverwriting: false, // Useless by design: there is only one ack waiting at a time
					AckTimeout:    0,     // Unused
					MaxRetries:    0,     // Unused
				},
				conn,
			)
			buf.ackCh = ackBuf.inCh
		}
	}

	for _, bufCfg := range outBufCfgs {
		if err := bufCfg.validate(); err != nil {
			return nil, err
		}
		buf := newOutBuffer(bufCfg, conn)
		b.outBuffers[bufCfg.ID] = buf
		if bufCfg.FrameType == arnetworkal.FrameTypeDataWithAck {
			// Automatically create an ack buffer...
			ackBufID := bufCfg.ID + ackBufferOffset
			// nolint: lll
			ackBuf := newInBuffer(
				InBufferConfig{
					ID:            ackBufID,
					FrameType:     arnetworkal.FrameTypeAck,
					Size:          1,     // Never more than one ack at a time
					MaxDataSize:   1,     // One byte of data: the sequence number
					IsOverwriting: false, // Useless by design: there is only one ack waiting at a time
				},
			)
			buf.ackCh = ackBuf.outCh
		}
	}

	go b.receiveFrames()

	return b, nil
}

func (b *bufferManager) receiveFrames() {
	for {
		netFrames, err := b.conn.Receive()
		if err != nil {
			log.Printf("error receiving frame: %s\n", err)
			continue
		}
		for _, netFrame := range netFrames {
			// Unpack the arnetworkal frame into an arnetwork frame...
			frame := Frame{
				seq:  netFrame.Seq,
				Data: netFrame.Data,
			}
			// Put the frame on the correct buffer...
			buf, ok := b.inBuffers[netFrame.ID]
			if !ok {
				// No buffer found to receive this frame.
				log.Printf(
					"received frame for unknown buffer %d\n",
					netFrame.ID,
				)
				continue
			}
			buf.inCh <- frame
		}
	}
}

func (b *bufferManager) OutC(bufID uint8) chan<- Frame {
	buf, ok := b.outBuffers[bufID]
	if !ok {
		return nil
	}
	return buf.inCh
}

func (b *bufferManager) InC(bufID uint8) <-chan Frame {
	buf, ok := b.inBuffers[bufID]
	if !ok {
		return nil
	}
	return buf.outCh
}
