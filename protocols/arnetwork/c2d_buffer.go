package arnetwork

import (
	"bytes"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type c2dBuffer struct {
	C2DBufferConfig
	buffer      *buffer
	inCh        chan Frame
	frameSender arnetworkal.FrameSender
	seq         uint8
	ackCh       chan Frame
}

func newC2DBuffer(
	bufCfg C2DBufferConfig,
	frameSender arnetworkal.FrameSender,
) *c2dBuffer {
	buf := &c2dBuffer{
		C2DBufferConfig: bufCfg,
		buffer:          newBuffer(bufCfg.ID, bufCfg.Size, bufCfg.IsOverwriting),
		frameSender:     frameSender,
	}

	log.WithField(
		"id",
		bufCfg.ID,
	).Debug("created new c2d frame buffer")

	go buf.receiveFrames()

	go buf.writeFrames()

	return buf
}

func (c *c2dBuffer) receiveFrames() {
	for frame := range c.inCh {
		if log.GetLevel() == log.DebugLevel {
			frame.uuid = uuid.NewV4().String()
		}
		c.buffer.inCh <- frame
	}
	close(c.buffer.inCh)
}

func (c *c2dBuffer) writeFrames() {
	log := log.WithField("id", c.buffer.id)
	log.Debug("c2d buffer is now buffering frames")
	for frame := range c.buffer.outCh {
		// Note that there's nothing we could do with the error here other than
		// log it, and we don't bother because writeFrame() already logs any
		// error that occurs since it's able to provide greater context than
		// this functions could-- i.e. how many attempt were made to deliver
		// the frame, etc.
		c.writeFrame(frame) // nolint: errcheck
	}
}

func (c *c2dBuffer) writeFrame(frame Frame) error { // nolint: unparam
	c.seq++ // Only increment seq once, no matter how many tries it takes
	netFrame := arnetworkal.Frame{
		UUID: frame.uuid,
		ID:   c.ID,
		Type: c.FrameType,
		Seq:  c.seq,
		Data: frame.Data,
	}
	log := log.WithField(
		"uuid",
		netFrame.UUID,
	).WithField(
		"id",
		c.buffer.id,
	).WithField(
		"seq",
		netFrame.Seq,
	)
	var attempts int
	for attempts = 0; attempts <= c.MaxRetries || c.MaxRetries == -1; attempts++ { // nolint: lll
		log.WithField(
			"attempt",
			attempts,
		).Debug("attempting to send arnetworkal frame")
		if err := c.frameSender.Send(netFrame); err != nil {
			log.WithField(
				"attempt",
				attempts,
			).Errorf("error sending arnetworkal frame: %s", err)
			return errors.Wrap(err, "error sending arnetworkal frame")
		}
		if netFrame.Type != arnetworkal.FrameTypeDataWithAck {
			return nil
		}
		select {
		case ack := <-c.ackCh:
			if bytes.Equal(
				[]byte(fmt.Sprintf("%d", netFrame.Seq)),
				ack.Data,
			) {
				return nil
			}
		case <-time.After(c.AckTimeout):
			log.WithField(
				"attempt",
				attempts,
			).Debug(
				"timed out waiting for acknowledgment of arnetworkal frame " +
					"receipt",
			)
		}
	}
	log.WithField(
		"attempts",
		attempts,
	).WithField(
		"retries",
		c.MaxRetries,
	).Error("exhausted retries sending arnetworkal frame")
	return errors.Errorf(
		"exhausted %d retries sending arnetworkal frame",
		c.MaxRetries,
	)
}
