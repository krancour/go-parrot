package arnetwork

import (
	log "github.com/Sirupsen/logrus"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

// NewBuffers returns maps of write-only channels for placing frames onto c2d
// buffers and read-only channels for receiving frames from d2c buffers. All
// channels are indexed by buffer ID.
func NewBuffers(
	frameSender arnetworkal.FrameSender,
	frameReceiver arnetworkal.FrameReceiver,
	c2dBufCfgs []C2DBufferConfig,
	d2cBufCfgs []D2CBufferConfig,
) (map[uint8]chan<- Frame, map[uint8]<-chan Frame, error) {
	c2dInChs := map[uint8]chan<- Frame{}
	d2cInChs := map[uint8]chan<- Frame{}
	d2cOutChs := map[uint8]<-chan Frame{}

	for _, bufCfg := range c2dBufCfgs {
		if err := bufCfg.validate(); err != nil {
			return nil, nil, err
		}
		buf := newC2DBuffer(bufCfg, frameSender)
		c2dInChs[bufCfg.ID] = buf.inCh
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
			return nil, nil, err
		}
		buf := newD2CBuffer(bufCfg)
		d2cInChs[bufCfg.ID] = buf.inCh
		d2cOutChs[bufCfg.ID] = buf.buffer.outCh
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
				frameSender,
			)
			buf.ackCh = ackBuf.buffer.inCh
		}
	}

	go receiveFrames(frameReceiver, d2cInChs)

	return c2dInChs, d2cOutChs, nil
}

func receiveFrames(
	frameReceiver arnetworkal.FrameReceiver,
	d2cInChs map[uint8]chan<- Frame,
) {
	for {
		netFrames, err := frameReceiver.Receive()
		if err != nil {
			log.Errorf("error receiving arnetworkal frames: %s", err)
			continue
		}
		for _, netFrame := range netFrames {
			// Find correct buffer for this frame...
			d2cInCh, ok := d2cInChs[netFrame.ID]
			if !ok {
				// No buffer found to receive this frame.
				log.WithField(
					"buffer",
					netFrame.ID,
				).Warn("received arnetworkal frame for unknown buffer")
				continue
			}
			// Unpack the arnetworkal frame into an arnetwork frame and put it
			// in the buffer...
			d2cInCh <- Frame{
				uuid: netFrame.UUID,
				seq:  netFrame.Seq,
				Data: netFrame.Data,
			}
		}
	}
}
