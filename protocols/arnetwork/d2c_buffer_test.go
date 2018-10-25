package arnetwork

import (
	"testing"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/stretchr/testify/require"
)

func TestD2CBufferReceiveFrames(t *testing.T) {
	testCases := []struct {
		name             string
		bufCfg           D2CBufferConfig
		initialBufRefSeq *uint8
		frame            Frame
		assertions       func(
			t *testing.T,
			expectedFrame Frame,
			frameCh <-chan Frame,
			ackCh <-chan Frame,
		)
	}{

		{
			name: "receive frame not needing acknowledgement and buffer has no " +
				"reference sequence number",
			bufCfg: D2CBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeData,
			},
			initialBufRefSeq: nil,
			frame: Frame{
				seq:  5,
				Data: []byte("foo"),
			},
			assertions: func(
				t *testing.T,
				expectedFrame Frame,
				frameCh <-chan Frame,
				_ <-chan Frame,
			) {
				_, ok := <-frameCh
				require.True(t, ok, "frame was not accepted, but should have been")
			},
		},

		{
			name: "receive frame not needing acknowledgement and frame sequence " +
				"number greater than buffer reference sequence number",
			bufCfg: D2CBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeData,
			},
			initialBufRefSeq: func() *uint8 { i := uint8(5); return &i }(),
			frame: Frame{
				seq:  10,
				Data: []byte("foo"),
			},
			assertions: func(
				t *testing.T,
				expectedFrame Frame,
				frameCh <-chan Frame,
				_ <-chan Frame,
			) {
				_, ok := <-frameCh
				require.True(t, ok, "frame was not accepted, but should have been")
			},
		},

		{
			name: "receive frame not needing acknowledgement and frame sequence " +
				"number less than buffer reference sequence number",
			bufCfg: D2CBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeData,
			},
			initialBufRefSeq: func() *uint8 { i := uint8(10); return &i }(),
			frame: Frame{
				seq:  5,
				Data: []byte("foo"),
			},
			assertions: func(
				t *testing.T,
				expectedFrame Frame,
				frameCh <-chan Frame,
				_ <-chan Frame,
			) {
				_, ok := <-frameCh
				require.False(t, ok, "frame was accepted and should not have been")
			},
		},

		{
			name: "receive frame not needing acknowledgement and frame sequence " +
				"number MUCH less than buffer reference sequence number",
			bufCfg: D2CBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeData,
			},
			initialBufRefSeq: func() *uint8 { i := uint8(25); return &i }(),
			frame: Frame{
				seq:  15,
				Data: []byte("foo"),
			},
			assertions: func(
				t *testing.T,
				expectedFrame Frame,
				frameCh <-chan Frame,
				_ <-chan Frame,
			) {
				_, ok := <-frameCh
				require.True(t, ok, "frame was not accepted, but should have been")
			},
		},

		{
			name: "receive frame requiring acknowledgement",
			bufCfg: D2CBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeDataWithAck,
			},
			initialBufRefSeq: func() *uint8 { i := uint8(5); return &i }(),
			frame: Frame{
				seq:  10,
				Data: []byte("foo"),
			},
			assertions: func(
				t *testing.T,
				expectedFrame Frame,
				frameCh <-chan Frame,
				ackCh <-chan Frame,
			) {
				_, ok := <-ackCh
				require.True(t, ok, "frame was not acknowledged, but should have been")
				_, ok = <-frameCh
				require.True(t, ok, "frame was not accepted, but should have been")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Remember that a new buffer will automatically begin receiving frame.
			// There's no need to explicitly call buf.receiveFrames()
			buf := newD2CBuffer(testCase.bufCfg)
			buf.seq = testCase.initialBufRefSeq
			// Set the channel that acks are written to so we can
			// make some assertions on it
			buf.ackCh = make(chan Frame)
			// Go ahead and give it a frame to deal with
			go func() {
				buf.inCh <- testCase.frame
				close(buf.inCh)
			}()
			// This isn't ideal, but because the buffer uses a goroutine to internally
			// copy from the input channel to the output channel, we need to wait for
			// a little while here to be sure that the frame that is accepted onto
			// the input channel has made it through to the output channel. If we
			// don't, we can't make reliable assertions about the state of the output
			// channel. We're expecting to find certain frame there, but it simply may
			// not be there yet.
			// <-time.After(5 * time.Second)
			testCase.assertions(t, testCase.frame, buf.buffer.outCh, buf.ackCh)
		})
	}
}
