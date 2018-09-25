package arnetwork

import (
	"log"

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
		buffer:          newBuffer(bufCfg.BaseBufferConfig),
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
	netFrame := arnetworkal.Frame{
		ID:   o.ID,
		Type: o.BaseBufferConfig.FrameType,
		Seq:  o.seq,
		Data: frame.Data,
	}
	o.seq++
	if err := o.conn.Send(netFrame); err != nil {
		log.Printf("error sending frame: %s\n", err)
	}
	if netFrame.Type == arnetworkal.FrameTypeDataWithAck {
		// TODO: Loop / wait for ack
	}
}
