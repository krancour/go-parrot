package arnetwork

import (
	"testing"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/krancour/go-parrot/protocols/arnetworkal/fake"

	"github.com/stretchr/testify/require"
)

func TestNewBufferManager(t *testing.T) {
	testCases := []struct {
		name       string
		c2dBufCfgs []C2DBufferConfig
		d2cBufCfgs []D2CBufferConfig
		assertions func(*testing.T, *bufferManager, error)
	}{

		{
			name: "invalid c2d buffer config",
			c2dBufCfgs: []C2DBufferConfig{
				{
					FrameType: arnetworkal.FrameTypeData,
					Size:      0,
				},
			},
			assertions: func(t *testing.T, _ *bufferManager, err error) {
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
			assertions: func(t *testing.T, _ *bufferManager, err error) {
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
			assertions: func(t *testing.T, bufMgr *bufferManager, err error) {
				require.NoError(t, err)

				require.Len(t, bufMgr.c2dBuffers, 1)
				c2dBuf, ok := bufMgr.c2dBuffers[5]
				require.True(t, ok)
				require.NotNil(t, c2dBuf.ackCh)

				require.Len(t, bufMgr.d2cBuffers, 1)
				d2cBuf, ok := bufMgr.d2cBuffers[10]
				require.True(t, ok)
				require.NotNil(t, d2cBuf.ackCh)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			bufMgrIface, err := NewBufferManager(
				&fake.FrameSender{},
				&fake.FrameReceiver{},
				testCase.c2dBufCfgs,
				testCase.d2cBufCfgs,
			)
			var bufMgr *bufferManager
			if err == nil {
				var ok bool
				bufMgr, ok = bufMgrIface.(*bufferManager)
				require.True(t, ok)
			}
			testCase.assertions(t, bufMgr, err)
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
	bufMgr := &bufferManager{
		frameReceiver: frameReceiver,
		d2cBuffers: map[uint8]*d2cBuffer{
			1: {
				inCh: testCh,
			},
		},
	}
	go bufMgr.receiveFrames()
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

func TestC2DCh(t *testing.T) {
	testCh := make(chan Frame)

	testCases := []struct {
		name       string
		bufMgr     *bufferManager
		bufID      uint8
		assertions func(*testing.T, chan<- Frame)
	}{

		{
			name:   "get non existing buffer",
			bufMgr: &bufferManager{},
			bufID:  10,
			assertions: func(t *testing.T, ch chan<- Frame) {
				require.Nil(t, ch)
			},
		},

		{
			name: "get existing buffer",
			bufMgr: &bufferManager{
				c2dBuffers: map[uint8]*c2dBuffer{
					10: {
						inCh: testCh,
					},
				},
			},
			bufID: 10,
			assertions: func(t *testing.T, ch chan<- Frame) {
				require.NotNil(t, ch)
				require.Equal(
					t,
					func(ch chan Frame) chan<- Frame {
						return ch
					}(testCh),
					ch,
				)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ch := testCase.bufMgr.C2DCh(testCase.bufID)
			testCase.assertions(t, ch)
		})
	}
}

func TestD2CCh(t *testing.T) {
	testCh := make(chan Frame)

	testCases := []struct {
		name       string
		bufMgr     *bufferManager
		bufID      uint8
		assertions func(*testing.T, <-chan Frame)
	}{

		{
			name:   "get non existing buffer",
			bufMgr: &bufferManager{},
			bufID:  10,
			assertions: func(t *testing.T, ch <-chan Frame) {
				require.Nil(t, ch)
			},
		},

		{
			name: "get existing buffer",
			bufMgr: &bufferManager{
				d2cBuffers: map[uint8]*d2cBuffer{
					10: {
						buffer: &buffer{
							outCh: testCh,
						},
					},
				},
			},
			bufID: 10,
			assertions: func(t *testing.T, ch <-chan Frame) {
				require.NotNil(t, ch)
				require.Equal(
					t,
					func(ch chan Frame) <-chan Frame {
						return ch
					}(testCh),
					ch,
				)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ch := testCase.bufMgr.D2CCh(testCase.bufID)
			testCase.assertions(t, ch)
		})
	}
}
