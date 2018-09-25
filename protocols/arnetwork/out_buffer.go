package arnetwork

import (
	"log"
	"sync"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

type outBuffer struct {
	OutBufferConfig
	*buffer
	conn arnetworkal.Connection
	mtx  sync.Mutex
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

// TODO: How do we deal with errors?
// TODO: How do we stop this goroutine?
func (o *outBuffer) writeFrames() {
	for {
		frame, ok, err := o.read()
		if err != nil {
			log.Printf(
				"error reading frame from head of outbound buffer: %s\n",
				err,
			)
			continue
		}
		if !ok {
			// We didn't get anything
			continue
		}
		o.writeFrame(frame)
	}
}

// TODO: How do we deal with errors?
func (o *outBuffer) writeFrame(frame Frame) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	netFrame := arnetworkal.Frame{
		ID:   frame.BufferID,
		Type: o.BaseBufferConfig.FrameType,
		Seq:  o.seq,
		Data: frame.Data,
	}
	o.seq++
	if err := o.conn.Send(netFrame); err != nil {
		log.Printf("error sending frame: %s\n", err)
	}
	if netFrame.Type == arnetworkal.FrameTypeDataWithAck {
		// TODO: Uncomment this
		// attempts := 1
		// TODO: We should have an ack coming back to us. How
		// do we learn about it?
		// TODO: Loop until we get an ack or reach max retries.
	}
}
