package arnetwork

// import (
// 	"container/ring"
// 	"testing"

// 	"github.com/stretchr/testify/require"
// )

// func TestBufferWrite(t *testing.T) {
// 	testCases := []struct {
// 		name           string
// 		bufCfg         BaseBufferConfig
// 		testFrames     []Frame
// 		expectedFrames []Frame
// 	}{
// 		{
// 			name: "fill buffer",
// 			bufCfg: BaseBufferConfig{
// 				Size: 5,
// 			},
// 			testFrames: []Frame{
// 				Frame{Data: []byte("a")},
// 				Frame{Data: []byte("b")},
// 				Frame{Data: []byte("c")},
// 				Frame{Data: []byte("d")},
// 				Frame{Data: []byte("e")},
// 			},
// 			expectedFrames: []Frame{
// 				Frame{Data: []byte("a")},
// 				Frame{Data: []byte("b")},
// 				Frame{Data: []byte("c")},
// 				Frame{Data: []byte("d")},
// 				Frame{Data: []byte("e")},
// 			},
// 		},
// 		{
// 			name: "overfill buffer with no overwrite",
// 			bufCfg: BaseBufferConfig{
// 				Size:          5,
// 				IsOverwriting: false,
// 			},
// 			testFrames: []Frame{
// 				Frame{Data: []byte("a")},
// 				Frame{Data: []byte("b")},
// 				Frame{Data: []byte("c")},
// 				Frame{Data: []byte("d")},
// 				Frame{Data: []byte("e")},
// 				Frame{Data: []byte("f")},
// 				Frame{Data: []byte("g")},
// 			},
// 			expectedFrames: []Frame{
// 				Frame{Data: []byte("a")},
// 				Frame{Data: []byte("b")},
// 				Frame{Data: []byte("c")},
// 				Frame{Data: []byte("d")},
// 				Frame{Data: []byte("e")},
// 			},
// 		},
// 		{
// 			name: "overfill buffer with overwrite",
// 			bufCfg: BaseBufferConfig{
// 				Size:          5,
// 				IsOverwriting: true,
// 			},
// 			testFrames: []Frame{
// 				Frame{Data: []byte("a")},
// 				Frame{Data: []byte("b")},
// 				Frame{Data: []byte("c")},
// 				Frame{Data: []byte("d")},
// 				Frame{Data: []byte("e")},
// 				Frame{Data: []byte("f")},
// 				Frame{Data: []byte("g")},
// 			},
// 			expectedFrames: []Frame{
// 				Frame{Data: []byte("c")},
// 				Frame{Data: []byte("d")},
// 				Frame{Data: []byte("e")},
// 				Frame{Data: []byte("f")},
// 				Frame{Data: []byte("g")},
// 			},
// 		},
// 	}
// 	for _, testCase := range testCases {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			buf := newBuffer(testCase.bufCfg)
// 			require.Nil(t, buf.head)
// 			require.NotNil(t, buf.next)
// 			for _, frame := range testCase.testFrames {
// 				buf.write(frame)
// 			}
// 			require.NotNil(t, buf.head)
// 			require.NotNil(t, buf.next)
// 			frames := []Frame{}
// 			buf.head.Do(func(frame interface{}) {
// 				frames = append(frames, frame.(Frame))
// 			})
// 			require.Equal(t, testCase.expectedFrames, frames)
// 		})
// 	}
// }

// func TestBufferReadFromEmptyBuffer(t *testing.T) {
// 	const bufSize = 5
// 	buf := newBuffer(
// 		BaseBufferConfig{
// 			Size: bufSize,
// 		},
// 	)
// 	require.Nil(t, buf.head)
// 	require.NotNil(t, buf.next)
// 	// Populate the ring buffer without using the write() function
// 	testFrames := []Frame{
// 		Frame{Data: []byte("a")},
// 		Frame{Data: []byte("b")},
// 		Frame{Data: []byte("c")},
// 		Frame{Data: []byte("d")},
// 		Frame{Data: []byte("e")},
// 	}
// 	buf.head = ring.New(bufSize)
// 	buf.next = buf.head
// 	for _, frame := range testFrames {
// 		buf.next.Value = frame
// 		buf.next = buf.next.Next()
// 	}
// 	for _, expectedFrame := range testFrames {
// 		frame, ok, err := buf.read()
// 		require.Nil(t, err)
// 		require.True(t, ok)
// 		require.Equal(t, expectedFrame, frame)
// 	}
// 	frame, ok, err := buf.read()
// 	require.Nil(t, err)
// 	require.False(t, ok)
// 	require.Equal(t, Frame{}, frame)
// }
