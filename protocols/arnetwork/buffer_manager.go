package arnetwork

import (
	"log"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

// BufferManager sends and receives arnetwork frames via a series of c2d and d2c
// buffers.
type BufferManager interface {
	C2DCh(uint8) chan<- Frame
	D2CCh(uint8) <-chan Frame
}

type bufferManager struct {
	conn       arnetworkal.Connection
	c2dBuffers map[uint8]*c2dBuffer
	d2cBuffers map[uint8]*d2cBuffer
}

// NewBufferManager returns a new BufferManager.
func NewBufferManager(
	conn arnetworkal.Connection,
	c2dBufCfgs []C2DBufferConfig,
	d2cBufCfgs []D2CBufferConfig,
) (BufferManager, error) {
	b := &bufferManager{
		conn: conn,
	}

	for _, bufCfg := range c2dBufCfgs {
		if err := bufCfg.validate(); err != nil {
			return nil, err
		}
		buf := newC2DBuffer(bufCfg, conn)
		b.c2dBuffers[bufCfg.ID] = buf
		if bufCfg.FrameType == arnetworkal.FrameTypeDataWithAck {
			// Automatically create an ack buffer...
			ackBufID := bufCfg.ID + ackBufferOffset
			// nolint: lll
			ackBuf := newD2CBuffer(
				D2CBufferConfig{
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

	for _, bufCfg := range d2cBufCfgs {
		if err := bufCfg.validate(); err != nil {
			return nil, err
		}
		buf := newD2CBuffer(bufCfg)
		b.d2cBuffers[bufCfg.ID] = buf
		if bufCfg.FrameType == arnetworkal.FrameTypeDataWithAck {
			// Automatically create an ack buffer...
			ackBufID := bufCfg.ID + ackBufferOffset
			// nolint: lll
			ackBuf := newC2DBuffer(
				C2DBufferConfig{
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
			buf, ok := b.d2cBuffers[netFrame.ID]
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

func (b *bufferManager) C2DCh(bufID uint8) chan<- Frame {
	buf, ok := b.c2dBuffers[bufID]
	if !ok {
		return nil
	}
	return buf.inCh
}

func (b *bufferManager) D2CCh(bufID uint8) <-chan Frame {
	buf, ok := b.d2cBuffers[bufID]
	if !ok {
		return nil
	}
	return buf.outCh
}
