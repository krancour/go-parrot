package arnetwork

import (
	log "github.com/Sirupsen/logrus"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
)

// NewBuffers returns maps of write-only channels for placing frames onto c2d
// buffers and read-only channels for receiving frames from d2c buffers. All
// channels are indexed by buffer ID.
func NewBuffers(
	frameSender arnetworkal.FrameSender,
	frameReceiver arnetworkal.FrameReceiver,
	c2dBufCfgs []C2DBufferConfig,
	d2cBufCfgs []D2CBufferConfig,
) (map[uint8]chan<- Frame, map[uint8]<-chan Frame, error) {
	c2dInChs := map[uint8]chan<- Frame{}
	d2cInChs := map[uint8]chan<- Frame{}
	d2cOutChs := map[uint8]<-chan Frame{}

	// TODO: This is a GUESS at how this buffer should be configured. The details
	// of this buffer are not well documented.
	pingBuf := newD2CBuffer(
		D2CBufferConfig{
			ID:            0,
			FrameType:     arnetworkal.FrameTypeData,
			Size:          20,
			MaxDataSize:   8, // This is the size of the "timespec" data we receive
			IsOverwriting: true,
		},
	)
	// This is for arnetwork internal use only. We'll add the input channel of
	// this d2cBuffer to d2cInChs so that receiveFrames(...) can add frames to
	// this buffer, BUT we'll refrain from adding the output channel of this
	// d2cBuffer to d2cOutChs, because those channels are returned from this
	// function and can be listened to by the caller-- which we do not want in
	// this case.
	d2cInChs[pingBuf.ID] = pingBuf.inCh

	// TODO: This is a GUESS at how this buffer should be configured. The
	// details of this buffer are not well documented.
	pongBuf := newC2DBuffer(
		C2DBufferConfig{
			ID:            1,
			FrameType:     arnetworkal.FrameTypeData,
			Size:          20,
			MaxDataSize:   8, // This is the size of the "timespec" data we must echo
			IsOverwriting: true,
		},
		frameSender,
	)
	// The above is for arnetwork internal use only. We won't add its input
	// channel to c2dInChs, because those channels are returned from this function
	// and can be written to by the caller-- which we do not want in this case.

	for _, bufCfg := range c2dBufCfgs {
		if err := bufCfg.validate(); err != nil {
			return nil, nil, err
		}
		buf := newC2DBuffer(bufCfg, frameSender)
		c2dInChs[bufCfg.ID] = buf.inCh
		if bufCfg.FrameType == arnetworkal.FrameTypeDataWithAck {
			// Automatically create an ack buffer...
			ackBufID := bufCfg.ID + ackBufferOffset
			// nolint: lll
			ackBuf := newD2CBuffer(
				D2CBufferConfig{
					ID:            ackBufID,
					FrameType:     arnetworkal.FrameTypeAck,
					Size:          1,     // Never more than one ack at a time
					MaxDataSize:   1,     // One byte of data: the sequence number
					IsOverwriting: false, // Useless by design: there is only one ack waiting at a time
				},
			)
			d2cInChs[ackBuf.ID] = ackBuf.inCh
			buf.ackCh = ackBuf.buffer.outCh
		}
	}

	for _, bufCfg := range d2cBufCfgs {
		if err := bufCfg.validate(); err != nil {
			return nil, nil, err
		}
		buf := newD2CBuffer(bufCfg)
		d2cInChs[bufCfg.ID] = buf.inCh
		d2cOutChs[bufCfg.ID] = buf.buffer.outCh
		if bufCfg.FrameType == arnetworkal.FrameTypeDataWithAck {
			// Automatically create an ack buffer...
			ackBufID := bufCfg.ID + ackBufferOffset
			// nolint: lll
			ackBuf := newC2DBuffer(
				C2DBufferConfig{
					ID:            ackBufID,
					FrameType:     arnetworkal.FrameTypeAck,
					Size:          1,     // Never more than one ack at a time
					MaxDataSize:   1,     // One byte of data: the sequence number
					IsOverwriting: false, // Useless by design: there is only one ack waiting at a time
					AckTimeout:    0,     // Unused
					MaxRetries:    0,     // Unused
				},
				frameSender,
			)
			buf.ackCh = ackBuf.buffer.inCh
		}
	}

	// Mux received frames into the appropriate buffers
	go receiveFrames(frameReceiver, d2cInChs)

	// Respond to pings. This turns out to be very important for avoiding
	// disconnects! Why? The arnetwork protocol (on the device end) assumes a
	// disconnect has occurred after five seconds of receiving no data from the
	// controller. Once that occurs, the device stops sending data as well. At
	// that point, neither controller or device is speaking to the other and the
	// disconnect is now real-- a self-fulfilling prophecy!
	//
	// Since we can't guarantee that a user of this SDK even HAS data
	// to send more often than every five seconds, our best bet for avoiding
	// disconnects is to simply respond to pings.
	go func() {
		for frame := range pingBuf.buffer.outCh {
			log.Debug("received ping; sending pong")
			pongBuf.inCh <- Frame{
				Data: frame.Data,
			}
		}
	}()

	return c2dInChs, d2cOutChs, nil
}

// receiveFrames muxes frames into the appropriate buffers.
func receiveFrames(
	frameReceiver arnetworkal.FrameReceiver,
	d2cInChs map[uint8]chan<- Frame,
) {
	for {
		netFrames, err := frameReceiver.Receive()
		if err != nil {
			log.Errorf("error receiving arnetworkal frames: %s", err)
			continue
		}
		for _, netFrame := range netFrames {
			// Find correct buffer for this frame...
			d2cInCh, ok := d2cInChs[netFrame.ID]
			if !ok {
				// No buffer found to receive this frame.
				log.WithField(
					"buffer",
					netFrame.ID,
				).Warn("received arnetworkal frame for unknown buffer")
				continue
			}
			// Unpack the arnetworkal frame into an arnetwork frame and put it
			// in the buffer...
			d2cInCh <- Frame{
				uuid: netFrame.UUID,
				seq:  netFrame.Seq,
				Data: netFrame.Data,
			}
		}
	}
}
