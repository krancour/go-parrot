package arnetwork

import (
	"fmt"
)

type inBuffer struct {
	InBufferConfig
	*buffer
	inCh  chan Frame
	seq   *uint8
	ackCh chan Frame
}

func newInBuffer(bufCfg InBufferConfig) *inBuffer {
	buf := &inBuffer{
		InBufferConfig: bufCfg,
		buffer:         newBuffer(bufCfg.BaseBufferConfig),
		inCh:           make(chan Frame),
	}

	go buf.receiveFrames()

	return buf
}

func (i *inBuffer) receiveFrames() {
	for frame := range i.inCh {
		// If acknowledgement was requested, send it...
		if i.ackCh != nil {
			i.ackCh <- Frame{
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
		if i.seq == nil || frame.seq > *i.seq || *i.seq-frame.seq >= 10 {
			i.seq = &frame.seq
			i.buffer.inCh <- frame
		}
	}
	close(i.buffer.inCh)
}
