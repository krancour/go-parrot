package arnetwork

import (
	log "github.com/Sirupsen/logrus"

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
		conn:       conn,
		c2dBuffers: map[uint8]*c2dBuffer{},
		d2cBuffers: map[uint8]*d2cBuffer{},
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
			buf.ackCh = ackBuf.buffer.outCh
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
			buf.ackCh = ackBuf.buffer.inCh
		}
	}

	go b.receiveFrames()

	return b, nil
}

func (b *bufferManager) receiveFrames() {
	for {
		netFrames, err := b.conn.Receive()
		if err != nil {
			log.Errorf("error receiving arnetworkal frames: %s", err)
			continue
		}
		for _, netFrame := range netFrames {
			// Find correct buffer for this frame...
			buf, ok := b.d2cBuffers[netFrame.ID]
			if !ok {
				// No buffer found to receive this frame.
				log.WithField(
					"buffer",
					netFrame.ID,
				).Error("received arnetworkal frame for unknown buffer")
				continue
			}
			// Unpack the arnetworkal frame into an arnetwork frame and put it
			// in the buffer...
			buf.inCh <- Frame{
				uuid: netFrame.UUID,
				seq:  netFrame.Seq,
				Data: netFrame.Data,
			}
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
	return buf.buffer.outCh
}
