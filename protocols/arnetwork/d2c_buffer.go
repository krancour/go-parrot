package arnetwork

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

type d2cBuffer struct {
	D2CBufferConfig
	buffer *buffer
	inCh   chan Frame
	seq    *uint8
	ackCh  chan Frame
}

func newD2CBuffer(bufCfg D2CBufferConfig) *d2cBuffer {
	buf := &d2cBuffer{
		D2CBufferConfig: bufCfg,
		buffer:          newBuffer(bufCfg.ID, bufCfg.Size, bufCfg.IsOverwriting),
		inCh:            make(chan Frame),
	}

	log.WithField(
		"id",
		bufCfg.ID,
	).Debug("created new d2c frame buffer")

	go buf.receiveFrames()

	return buf
}

func (d *d2cBuffer) receiveFrames() {
	log := log.WithField("id", d.buffer.id)
	for frame := range d.inCh {
		log = log.WithField("uuid", frame.uuid).WithField("seq", frame.seq)
		// If acknowledgement was requested, send it...
		if d.FrameType == arnetworkal.FrameTypeDataWithAck && d.ackCh != nil {
			log.Debug("acknowledging receipt of frame")
			d.ackCh <- Frame{
				Data: []byte(fmt.Sprintf("%d", frame.seq)),
			}
		}
		// If this buffer doesn't yet have a reference sequence number or...
		//
		// The frame's sequence number is higher than the buffer's reference
		// sequence number or...
		//
		// The gap between the frame sequence number and the buffer's reference
		// sequence number is bigger than 10...
		//
		// Then we should accept this frame and update the buffer's reference
		// sequence number.
		//
		// Otherwise, we're dealing with a frame that has arrived out of order.
		// We'll drop such frames.
		if d.seq == nil || frame.seq > *d.seq || *d.seq-frame.seq >= 10 {
			log.Debug("accepting frame")
			d.seq = &frame.seq
			d.buffer.inCh <- frame
		} else {
			log.Debug("frame appears to be a duplicate or out of sequence; " +
				"dropping it")
		}
	}
	if d.ackCh != nil {
		close(d.ackCh)
	}
	close(d.buffer.inCh)
}
