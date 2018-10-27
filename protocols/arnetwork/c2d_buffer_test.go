package arnetwork

import (
	"fmt"
	"testing"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/krancour/go-parrot/protocols/arnetworkal/fake"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestWriteFrame(t *testing.T) {
	testCases := []struct {
		name         string
		bufCfg       C2DBufferConfig
		sendBehavior func(
			netFrame arnetworkal.Frame,
			callCount int,
			ackCh chan<- Frame,
		) error
		assertions func(t *testing.T, err error, sendCallCount int)
	}{

		{
			name: "write frame with no ack required",
			bufCfg: C2DBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeData,
				Size:      1,
			},
			sendBehavior: func(arnetworkal.Frame, int, chan<- Frame) error {
				// Simulate no error sending
				return nil
			},
			assertions: func(t *testing.T, err error, sendCallCount int) {
				require.NoError(t, err)
				require.Equal(t, 1, sendCallCount)
			},
		},

		{
			name: "try to write frame and get arnetworkal error",
			bufCfg: C2DBufferConfig{
				ID:        1,
				FrameType: arnetworkal.FrameTypeData,
				Size:      1,
			},
			sendBehavior: func(arnetworkal.Frame, int, chan<- Frame) error {
				return errors.New("error sending arnetworkal frame")
			},
			assertions: func(t *testing.T, err error, sendCallCount int) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "arnetworkal")
				require.Equal(t, 1, sendCallCount)
			},
		},

		{
			name: "write frame with ack required",
			bufCfg: C2DBufferConfig{
				ID:         1,
				FrameType:  arnetworkal.FrameTypeDataWithAck,
				Size:       1,
				AckTimeout: time.Second,
			},
			sendBehavior: func(
				netFrame arnetworkal.Frame,
				_ int,
				ackCh chan<- Frame,
			) error {
				go func() {
					ackCh <- Frame{
						Data: []byte(fmt.Sprintf("%d", netFrame.Seq)),
					}
				}()
				return nil
			},
			assertions: func(t *testing.T, err error, sendCallCount int) {
				require.NoError(t, err)
				require.Equal(t, 1, sendCallCount)
			},
		},

		{
			name: "failure to write frame with ack followed by successful retry",
			bufCfg: C2DBufferConfig{
				ID:         1,
				FrameType:  arnetworkal.FrameTypeDataWithAck,
				Size:       1,
				AckTimeout: time.Second,
				MaxRetries: 1,
			},
			sendBehavior: func(
				netFrame arnetworkal.Frame,
				callCount int,
				ackCh chan<- Frame,
			) error {
				if callCount == 1 {
					// Simulate no error sending-- but no ack of receipt from the device
					return nil
				}
				go func() {
					ackCh <- Frame{
						Data: []byte(fmt.Sprintf("%d", netFrame.Seq)),
					}
				}()
				return nil
			},
			assertions: func(t *testing.T, err error, sendCallCount int) {
				require.NoError(t, err)
				require.Equal(t, 2, sendCallCount)
			},
		},

		{
			name: "failure to write frame with ack and exhausted retries",
			bufCfg: C2DBufferConfig{
				ID:         1,
				FrameType:  arnetworkal.FrameTypeDataWithAck,
				Size:       1,
				AckTimeout: time.Second,
				MaxRetries: 1,
			},
			sendBehavior: func(arnetworkal.Frame, int, chan<- Frame) error {
				// Simulate no error sending-- but no ack of receipt from the device
				return nil
			},
			assertions: func(t *testing.T, err error, sendCallCount int) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "exhausted")
				require.Equal(t, 2, sendCallCount)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			conn := &fake.Connection{}
			err := testCase.bufCfg.validate()
			require.NoError(t, err)
			buf := newC2DBuffer(testCase.bufCfg, conn)
			buf.ackCh = make(chan Frame)
			initialSeq := buf.seq
			frame := Frame{
				Data: []byte("foo"),
			}
			sendCallCount := 0
			conn.SendBehavior = func(netFrame arnetworkal.Frame) error {
				sendCallCount++
				require.Equal(t, buf.ID, netFrame.ID)
				require.Equal(t, buf.FrameType, netFrame.Type)
				require.Equal(t, buf.seq, netFrame.Seq)
				require.Equal(t, frame.Data, netFrame.Data)
				return testCase.sendBehavior(netFrame, sendCallCount, buf.ackCh)
			}
			err = buf.writeFrame(frame)
			require.Equal(t, initialSeq+1, buf.seq)
			testCase.assertions(t, err, sendCallCount)
		})
	}
}
