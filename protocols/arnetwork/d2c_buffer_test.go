package arnetwork

import (
	"fmt"
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
			buf *d2cBuffer,
		)
	}{

		{
			name: "receive frame not needing acknowledgement and buffer has no " +
				"reference sequence number",
			bufCfg: D2CBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeData,
				Size:      1,
			},
			initialBufRefSeq: nil,
			frame: Frame{
				seq:  5,
				Data: []byte("foo"),
			},
			assertions: func(
				t *testing.T,
				expectedFrame Frame,
				buf *d2cBuffer,
			) {
				_, ok := <-buf.outCh
				require.True(t, ok, "frame was not accepted, but should have been")
				require.Equal(t, expectedFrame.seq, *buf.seq)
			},
		},

		{
			name: "receive frame not needing acknowledgement and frame sequence " +
				"number greater than buffer reference sequence number",
			bufCfg: D2CBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeData,
				Size:      1,
			},
			initialBufRefSeq: func() *uint8 { i := uint8(5); return &i }(),
			frame: Frame{
				seq:  10,
				Data: []byte("foo"),
			},
			assertions: func(
				t *testing.T,
				expectedFrame Frame,
				buf *d2cBuffer,
			) {
				_, ok := <-buf.outCh
				require.True(t, ok, "frame was not accepted, but should have been")
				require.Equal(t, expectedFrame.seq, *buf.seq)
			},
		},

		{
			name: "receive frame not needing acknowledgement and frame sequence " +
				"number less than buffer reference sequence number",
			bufCfg: D2CBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeData,
				Size:      1,
			},
			initialBufRefSeq: func() *uint8 { i := uint8(10); return &i }(),
			frame: Frame{
				seq:  5,
				Data: []byte("foo"),
			},
			assertions: func(
				t *testing.T,
				expectedFrame Frame,
				buf *d2cBuffer,
			) {
				_, ok := <-buf.outCh
				require.False(t, ok, "frame was accepted and should not have been")
				require.Equal(t, uint8(10), *buf.seq)
			},
		},

		{
			name: "receive frame not needing acknowledgement and frame sequence " +
				"number equal to buffer reference sequence number",
			bufCfg: D2CBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeData,
				Size:      1,
			},
			initialBufRefSeq: func() *uint8 { i := uint8(10); return &i }(),
			frame: Frame{
				seq:  10,
				Data: []byte("foo"),
			},
			assertions: func(
				t *testing.T,
				expectedFrame Frame,
				buf *d2cBuffer,
			) {
				_, ok := <-buf.outCh
				require.False(t, ok, "frame was accepted and should not have been")
				require.Equal(t, uint8(10), *buf.seq)
			},
		},

		{
			name: "receive frame not needing acknowledgement and frame sequence " +
				"number MUCH less than buffer reference sequence number",
			bufCfg: D2CBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeData,
				Size:      1,
			},
			initialBufRefSeq: func() *uint8 { i := uint8(25); return &i }(),
			frame: Frame{
				seq:  15,
				Data: []byte("foo"),
			},
			assertions: func(
				t *testing.T,
				expectedFrame Frame,
				buf *d2cBuffer,
			) {
				_, ok := <-buf.outCh
				require.True(t, ok, "frame was not accepted, but should have been")
				require.Equal(t, expectedFrame.seq, *buf.seq)
			},
		},

		{
			name: "receive frame requiring acknowledgement",
			bufCfg: D2CBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeDataWithAck,
				Size:      1,
			},
			initialBufRefSeq: func() *uint8 { i := uint8(5); return &i }(),
			frame: Frame{
				seq:  10,
				Data: []byte("foo"),
			},
			assertions: func(
				t *testing.T,
				expectedFrame Frame,
				buf *d2cBuffer,
			) {
				ackFrame, ok := <-buf.ackCh
				require.True(t, ok, "frame was not acknowledged, but should have been")
				require.Equal(
					t,
					[]byte(fmt.Sprintf("%d", expectedFrame.seq)),
					ackFrame.Data,
				)
				_, ok = <-buf.outCh
				require.True(t, ok, "frame was not accepted, but should have been")
				require.Equal(t, expectedFrame.seq, *buf.seq)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.bufCfg.validate()
			require.NoError(t, err)
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
			testCase.assertions(t, testCase.frame, buf)
		})
	}
}
