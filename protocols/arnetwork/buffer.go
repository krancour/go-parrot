package arnetwork

import (
	"container/ring"
	"fmt"
	"sync"
)

const ackBufferOffset uint8 = 128

// buffer implements a ring buffer. Behavior when writing to a full buffer is
// determined by BufferConfig.IsOverwriting. Outbound buffers continuously
// read from their own head to send frames out via the arnetworkal connection,
// while inbound buffers continuously read from their own head to invoke the
// specified function.
// TODO: Can buffer be improved to use channels instead of locks?
// https://content.pivotal.io/blog/a-channel-based-ring-buffer-in-go
type buffer struct {
	mtx           sync.Mutex
	head          *ring.Ring
	next          *ring.Ring
	isOverwriting bool
	seq           uint8
}

func newBuffer(bufCfg BaseBufferConfig) *buffer {
	return &buffer{
		next:          ring.New(int(bufCfg.Size)),
		isOverwriting: bufCfg.IsOverwriting,
	}
}

// write attempts to add a frame to a buffer and returns a boolean indicating
// success or failure.
func (b *buffer) write(frame Frame) bool {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if b.next == b.head {
		// The next place to write to is the head. i.e. The buffer is full.
		// How do we deal with it?
		// Are we "overwriting?"
		if !b.isOverwriting {
			// We're not "overwriting," so we'll just drop the new frame by
			// simply returning.
			return false
		}
		// We are "overwriting," so we'll advance the head (i.e. drop oldest
		// frame in the buffer.
		b.head = b.head.Next()
	}
	b.next.Value = frame
	if b.head == nil {
		b.head = b.next
	}
	b.next = b.next.Next()
	return true
}

// read reads a single frame from the head of the buffer. A boolean is also
// returned to indicate success or failure. It is possible, for instance, for a
// read ot return nothing without the occurence of an error, simply because the
// buffer is empty. An error is also potentially returned.
func (b *buffer) read() (Frame, bool, error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	if b.head == nil {
		// The buffer is empty.
		return Frame{}, false, nil
	}
	frame, ok := b.head.Value.(Frame)
	b.head = b.head.Next()
	// If we've advanced the head to the next spot we will write to, we have an
	// empty buffer and we should represent that by pointing the head to nil.
	if b.head == b.next {
		b.head = nil
	}
	if !ok {
		return frame,
			false,
			fmt.Errorf("data read from buffer was not an arnetwork frame")
	}
	return frame, true, nil
}
