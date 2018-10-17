package arnetwork

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

type outBuffer struct {
	OutBufferConfig
	*buffer
	conn  arnetworkal.Connection
	seq   uint8
	ackCh chan Frame
}

func newOutBuffer(
	bufCfg OutBufferConfig,
	conn arnetworkal.Connection,
) *outBuffer {
	buf := &outBuffer{
		OutBufferConfig: bufCfg,
		buffer:          newBuffer(bufCfg.Size, bufCfg.IsOverwriting),
		conn:            conn,
	}

	go buf.writeFrames()

	return buf
}

func (o *outBuffer) writeFrames() {
	for frame := range o.outCh {
		o.writeFrame(frame)
	}
}

func (o *outBuffer) writeFrame(frame Frame) {
	for attempts := 0; attempts <= o.MaxRetries || o.MaxRetries == -1; attempts++ {
		netFrame := arnetworkal.Frame{
			ID:   o.ID,
			Type: o.FrameType,
			Seq:  o.seq,
			Data: frame.Data,
		}
		o.seq++
		if err := o.conn.Send(netFrame); err != nil {
			log.Printf("error sending frame: %s\n", err)
		}
		if netFrame.Type == arnetworkal.FrameTypeDataWithAck {
			select {
			case ack := <-o.ackCh:
				if bytes.Equal(
					[]byte(fmt.Sprintf("%d", netFrame.Seq)),
					ack.Data,
				) {
					return
				}
			case <-time.After(o.AckTimeout):
			}
		}
	}
}
