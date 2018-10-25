package arnetwork

import (
	"bytes"
	"fmt"
	"log"
	"time"

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
		buffer:          newBuffer(bufCfg.Size, bufCfg.IsOverwriting),
		conn:            conn,
	}

	go buf.writeFrames()

	return buf
}

func (c *c2dBuffer) writeFrames() {
	for frame := range c.outCh {
		if err := c.writeFrame(frame); err != nil {
			log.Println(err)
		}
	}
}

func (c *c2dBuffer) writeFrame(frame Frame) error {
	c.seq++ // Only increment seq once, no matter how many tries it takes
	for attempts := 0; attempts <= c.MaxRetries || c.MaxRetries == -1; attempts++ {
		log.Printf("attempt %d to send frame", attempts)
		netFrame := arnetworkal.Frame{
			ID:   c.ID,
			Type: c.FrameType,
			Seq:  c.seq,
			Data: frame.Data,
		}
		if err := c.conn.Send(netFrame); err != nil {
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
			log.Println("timed out waiting for acknowledgment of frame receipt")
		}
	}
	return fmt.Errorf("exhausted %d retries sending frame", c.MaxRetries)
}
