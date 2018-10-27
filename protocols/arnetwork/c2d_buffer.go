package arnetwork

import (
	"bytes"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

type c2dBuffer struct {
	C2DBufferConfig
	*buffer
	conn  arnetworkal.Connection
	seq   uint8
	ackCh chan Frame
}

func newC2DBuffer(
	bufCfg C2DBufferConfig,
	conn arnetworkal.Connection,
) *c2dBuffer {
	buf := &c2dBuffer{
		C2DBufferConfig: bufCfg,
		buffer:          newBuffer(bufCfg.ID, bufCfg.Size, bufCfg.IsOverwriting),
		conn:            conn,
	}

	log.WithField(
		"id",
		bufCfg.ID,
	).Debug("created new c2d frame buffer")

	go buf.writeFrames()

	return buf
}

func (c *c2dBuffer) writeFrames() {
	log := log.WithField("id", c.id)
	log.Debug("c2d buffer is now buffering frames")
	for frame := range c.outCh {
		c.writeFrame(frame) // nolint: errcheck
	}
}

func (c *c2dBuffer) writeFrame(frame Frame) error { // nolint: unparam
	c.seq++ // Only increment seq once, no matter how many tries it takes
	netFrame := arnetworkal.Frame{
		ID:   c.ID,
		Type: c.FrameType,
		Seq:  c.seq,
		Data: frame.Data,
	}
	log := log.WithField(
		"id",
		c.id,
	).WithField(
		"seq",
		netFrame.Seq,
	)
	var attempts int
	for attempts = 0; attempts <= c.MaxRetries || c.MaxRetries == -1; attempts++ {
		log.WithField(
			"attempt",
			attempts,
		).Debug("attempting to send frame")
		if err := c.conn.Send(netFrame); err != nil {
			log.WithField(
				"attempt",
				attempts,
			).Errorf("error sending frame: %s", err)
			return fmt.Errorf("error sending frame: %s", err)
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
			).Debug("timed out waiting for acknowledgment of frame receipt")
		}
	}
	log.WithField(
		"attempts",
		attempts,
	).WithField(
		"retries",
		c.MaxRetries,
	).Error("exhausted retries sending frame")
	return fmt.Errorf("exhausted %d retries sending frame", c.MaxRetries)
}
