package arnetwork

import log "github.com/Sirupsen/logrus"

const ackBufferOffset uint8 = 128

type buffer struct {
	id            uint8
	inCh          chan Frame
	outCh         chan Frame
	isOverwriting bool
}

func newBuffer(id uint8, size int32, isOverwriting bool) *buffer {
	buf := &buffer{
		id:            id,
		inCh:          make(chan Frame),
		outCh:         make(chan Frame, size),
		isOverwriting: isOverwriting,
	}

	log.WithField(
		"id",
		id,
	).WithField(
		"overwriting",
		isOverwriting,
	).Debug("created new arnetwork frame buffer")

	go buf.bufferFrames()

	return buf
}

func (b *buffer) bufferFrames() {
	log := log.WithField("id", b.id)
	log.Debug("buffer is now buffering arnetwork frames")
	for frame := range b.inCh {
		select {
		// Try to write the frame to the outCh.
		case b.outCh <- frame:
			// Success! We're done.
			log.Debug("buffer had room for new arnetwork frame")
		default: // outCh is full...
			// If we're not "overwriting," just drop the frame and continue on
			// to the next one.
			if !b.isOverwriting {
				log.Debug(
					"buffer is full and overwriting is not enabled; dropping new " +
						"arnetwork frame",
				)
				continue
			}
			select {
			// Try to remove the oldest frame from the outCh.
			case <-b.outCh:
				// Success! We made room on the outCh.
				log.Debug(
					"buffer is full and overwriting is enabled; dropping oldest " +
						"arnetwork frame",
				)
			default:
				// outCh is already empty. Good!
				log.Debug("buffer had room for new arnetwork frame")
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
