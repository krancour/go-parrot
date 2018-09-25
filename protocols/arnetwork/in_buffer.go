package arnetwork

import "sync"

type inBuffer struct {
	InBufferConfig
	*buffer
	mtx            sync.Mutex
	seqInitialized bool
}

func newInBuffer(bufCfg InBufferConfig) *inBuffer {
	buf := &inBuffer{
		InBufferConfig: bufCfg,
		buffer:         newBuffer(bufCfg.BaseBufferConfig),
	}
	// TODO: Start a goroutine that reads from the head of the buffer and passes
	// the frame to a function.
	return buf
}

func (i *inBuffer) write(frame Frame) bool {
	i.mtx.Lock()
	defer i.mtx.Unlock()
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
	// We won't write it to the underlying buffer, but we'll still return true,
	// which will allow acknowledgement of this frame to proceed regardless.
	if !i.seqInitialized ||
		frame.seq > i.seq ||
		i.seq-frame.seq >= 10 {
		i.seq = frame.seq
		i.seqInitialized = true
	} else {
		return true
	}
	return i.buffer.write(frame)
}
