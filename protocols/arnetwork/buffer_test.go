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
			buf := newBuffer(testCase.size, testCase.isOverwriting)
			for _, frame := range testCase.testFrames {
				select {
				case buf.inCh <- frame:
				case <-time.After(time.Second):
					require.Fail(t, "write to buffer unexpectedly timed out")
				}
			}
			// This isn't ideal, but because the buffer uses a goroutine to internally
			// copy from the input channel to the output channel, we need to wait for
			// a little while here to be sure that everything we put on the input
			// channel has made it through to the output channel. If we don't, we
			// can't make reliable assertions about the things on the output channel.
			// We're expecting to find certain frames there, but they simply may not
			// be there yet.
			<-time.After(2 * time.Second)
			frames := []Frame{}
		loop:
			for {
				select {
				case frame := <-buf.outCh:
					frames = append(frames, frame)
				default:
					break loop
				}
			}
			require.Equal(t, testCase.expectedFrames, frames)
		})
	}
}

func TestEmptyBuffer(t *testing.T) {
	const bufSize = 5
	buf := newBuffer(bufSize, false)
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
