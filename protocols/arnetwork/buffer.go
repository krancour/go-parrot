package arnetwork

const ackBufferOffset uint8 = 128

type buffer struct {
	inCh          chan Frame
	outCh         chan Frame
	isOverwriting bool
}

func newBuffer(bufCfg BaseBufferConfig) *buffer {
	buf := &buffer{
		inCh:          make(chan Frame),
		outCh:         make(chan Frame, bufCfg.Size),
		isOverwriting: bufCfg.IsOverwriting,
	}

	go buf.receiveFrames()

	return buf
}

func (b *buffer) receiveFrames() {
	for frame := range b.inCh {
		select {
		// Try to write the frame to the outCh.
		case b.outCh <- frame: // Success! We're done.
		default: // outCh is full...
			// If we're not "overwriting," just drop the frame and continue on
			// to the next one.
			if !b.isOverwriting {
				continue
			}
			select {
			// Try to remove the oldest frame from the outCh.
			case <-b.outCh: // Success! We made room on the outCh.
			default: // outCh is already empty. Good!
			}
			// We make an assumption here that nobody writes to outCh besides
			// this function. This is a relatively safe assumption because
			// nothing else in this package does and the outCh is NOT exported.
			// Given this assumption, after the lines above, we're guaranteed to
			// have space on the outCh.
			b.outCh <- frame
		}
	}
	// Signal anyone listening to the outCh that there's no more data coming!
	close(b.outCh)
}
