package arnetwork

import (
	"fmt"
)

type d2cBuffer struct {
	D2CBufferConfig
	*buffer
	inCh  chan Frame
	seq   *uint8
	ackCh chan Frame
}

func newD2CBuffer(bufCfg D2CBufferConfig) *d2cBuffer {
	buf := &d2cBuffer{
		D2CBufferConfig: bufCfg,
		buffer:          newBuffer(bufCfg.Size, bufCfg.IsOverwriting),
		inCh:            make(chan Frame),
	}

	go buf.receiveFrames()

	return buf
}

func (d *d2cBuffer) receiveFrames() {
	for frame := range d.inCh {
		// If acknowledgement was requested, send it...
		if d.ackCh != nil {
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
		// sequecne number is bigger than 10...
		//
		// Then we should accept this frame and update the buffer's reference
		// sequence number.
		//
		// Otherwise, we're dealing with a frame that has arrived out of order.
		// We'll drop such frames.
		if d.seq == nil || frame.seq > *d.seq || *d.seq-frame.seq >= 10 {
			d.seq = &frame.seq
			d.buffer.inCh <- frame
		}
	}
	close(d.buffer.inCh)
}
