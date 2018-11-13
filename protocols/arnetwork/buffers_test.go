package arnetwork

import (
	"testing"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/krancour/go-parrot/protocols/arnetworkal/fake"

	"github.com/stretchr/testify/require"
)

func TestNewBuffers(t *testing.T) {
	testCases := []struct {
		name       string
		c2dBufCfgs []C2DBufferConfig
		d2cBufCfgs []D2CBufferConfig
		assertions func(
			*testing.T,
			map[uint8]chan<- Frame,
			map[uint8]<-chan Frame,
			error,
		)
	}{

		{
			name: "invalid c2d buffer config",
			c2dBufCfgs: []C2DBufferConfig{
				{
					FrameType: arnetworkal.FrameTypeData,
					Size:      0,
				},
			},
			assertions: func(
				t *testing.T,
				_ map[uint8]chan<- Frame,
				_ map[uint8]<-chan Frame,
				err error,
			) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid")
			},
		},

		{
			name: "invalid d2c buffer config",
			d2cBufCfgs: []D2CBufferConfig{
				{
					FrameType: arnetworkal.FrameTypeData,
					Size:      0,
				},
			},
			assertions: func(
				t *testing.T,
				_ map[uint8]chan<- Frame,
				_ map[uint8]<-chan Frame,
				err error,
			) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid")
			},
		},

		{
			name: "valid configuration",
			c2dBufCfgs: []C2DBufferConfig{
				{
					ID:        5,
					FrameType: arnetworkal.FrameTypeDataWithAck,
					Size:      1,
				},
			},
			d2cBufCfgs: []D2CBufferConfig{
				{
					ID:        10,
					FrameType: arnetworkal.FrameTypeDataWithAck,
					Size:      1,
				},
			},
			assertions: func(
				t *testing.T,
				c2dChs map[uint8]chan<- Frame,
				d2cChs map[uint8]<-chan Frame,
				err error,
			) {
				require.NoError(t, err)

				require.Len(t, c2dChs, 1)
				_, ok := c2dChs[5]
				require.True(t, ok)

				require.Len(t, d2cChs, 1)
				_, ok = d2cChs[10]
				require.True(t, ok)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c2dChs, d2cChs, err := NewBuffers(
				&fake.FrameSender{},
				&fake.FrameReceiver{},
				testCase.c2dBufCfgs,
				testCase.d2cBufCfgs,
			)
			testCase.assertions(t, c2dChs, d2cChs, err)
		})
	}
}

func TestReceiveFrames(t *testing.T) {
	const numFrames = 5
	callCount := 0
	testFrame := arnetworkal.Frame{
		ID:   1,
		Seq:  2,
		Data: []byte("foo"),
	}
	frameReceiver := &fake.FrameReceiver{
		ReceiveBehavior: func() ([]arnetworkal.Frame, error) {
			if callCount > 0 {
				return nil, nil
			}
			callCount++
			frames := []arnetworkal.Frame{}
			for i := 0; i < numFrames; i++ {
				frames = append(frames, testFrame)
			}
			return frames, nil
		},
	}
	testCh := make(chan Frame)
	go receiveFrames(frameReceiver, map[uint8]chan<- Frame{1: testCh})
	for i := 0; i < numFrames; i++ {
		select {
		case frame, ok := <-testCh:
			require.True(t, ok)
			require.Equal(t, testFrame.Seq, frame.seq)
			require.Equal(t, testFrame.Data, frame.Data)
		case <-time.After(time.Second):
			require.Fail(t, "timed out waiting for frame")
		}
	}
	select {
	case <-testCh:
		require.Fail(t, "received frame, but should not have")
	case <-time.After(2 * time.Second):
	}
}
