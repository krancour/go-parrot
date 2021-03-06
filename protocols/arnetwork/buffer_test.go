package arnetwork

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBuffer(t *testing.T) {
	testCases := []struct {
		name           string
		size           int32
		isOverwriting  bool
		testFrames     []Frame
		expectedFrames []Frame
	}{

		{
			name:          "fill buffer",
			size:          5,
			isOverwriting: false,
			testFrames: []Frame{
				{Data: []byte("a")},
				{Data: []byte("b")},
				{Data: []byte("c")},
				{Data: []byte("d")},
				{Data: []byte("e")},
			},
			expectedFrames: []Frame{
				{Data: []byte("a")},
				{Data: []byte("b")},
				{Data: []byte("c")},
				{Data: []byte("d")},
				{Data: []byte("e")},
			},
		},

		{
			name:          "overfill buffer with no overwrite",
			size:          5,
			isOverwriting: false,
			testFrames: []Frame{
				{Data: []byte("a")},
				{Data: []byte("b")},
				{Data: []byte("c")},
				{Data: []byte("d")},
				{Data: []byte("e")},
				{Data: []byte("f")},
				{Data: []byte("g")},
			},
			expectedFrames: []Frame{
				{Data: []byte("a")},
				{Data: []byte("b")},
				{Data: []byte("c")},
				{Data: []byte("d")},
				{Data: []byte("e")},
			},
		},

		{
			name:          "overfill buffer with overwrite",
			size:          5,
			isOverwriting: true,
			testFrames: []Frame{
				{Data: []byte("a")},
				{Data: []byte("b")},
				{Data: []byte("c")},
				{Data: []byte("d")},
				{Data: []byte("e")},
				{Data: []byte("f")},
				{Data: []byte("g")},
			},
			expectedFrames: []Frame{
				{Data: []byte("c")},
				{Data: []byte("d")},
				{Data: []byte("e")},
				{Data: []byte("f")},
				{Data: []byte("g")},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			buf := newBuffer(1, testCase.size, testCase.isOverwriting)
			for _, frame := range testCase.testFrames {
				select {
				case buf.inCh <- frame:
				case <-time.After(time.Second):
					require.Fail(t, "write to buffer unexpectedly timed out")
				}
			}
			close(buf.inCh)
			select {
			case <-buf.doneCh:
			case <-time.After(2 * time.Second):
				require.Fail(t, "timed out waiting for buffer to fill")
			}
			frames := []Frame{}
		loop:
			for {
				select {
				case frame, ok := <-buf.outCh:
					if !ok {
						break loop
					}
					frames = append(frames, frame)
				case <-time.After(time.Second):
					require.Fail(t, "read from buffer unexpectedly timed out")
				}
			}
			require.Equal(t, testCase.expectedFrames, frames)
		})
	}
}

func TestEmptyBuffer(t *testing.T) {
	const bufSize = 5
	buf := newBuffer(1, bufSize, false)
	testFrames := []Frame{
		{Data: []byte("a")},
		{Data: []byte("b")},
		{Data: []byte("c")},
		{Data: []byte("d")},
		{Data: []byte("e")},
	}
	for _, frame := range testFrames {
		select {
		case buf.inCh <- frame:
		case <-time.After(time.Second):
			require.Fail(t, "write to buffer unexpectedly timed out")
		}
	}
	for _, expectedFrame := range testFrames {
		select {
		case frame := <-buf.outCh:
			require.Equal(t, expectedFrame, frame)
		case <-time.After(time.Second):
			require.Fail(t, "read from buffer unexpectedly timed out")
		}
	}
	// The next read from the buffer SHOULD be blocked
	select {
	case <-buf.outCh:
		require.Fail(
			t,
			"unexpectedly succeeded in reading past the end of the buffer",
		)
	default:
	}
}
