package wifi

import (
	"testing"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/stretchr/testify/assert"
)

func TestEncodeFrame(t *testing.T) {
	packet := defaultEncodeFrame(
		arnetworkal.Frame{
			Type: arnetworkal.FrameTypeAck,
			ID:   186,
			Seq:  39,
			Data: []byte{0x42},
		},
	)
	assert.Equal(
		t,
		[]byte{
			0x01,                   // Type
			0xba,                   // ID
			0x27,                   // Seq
			0x08, 0x00, 0x00, 0x00, // Frame size (little endian)
			0x42, // Data
		},
		packet,
	)
}

func TestDecodePacket(t *testing.T) {
	testCases := []struct {
		name   string
		packet []byte
		assert func(t *testing.T, frames []arnetworkal.Frame, err error)
	}{
		{
			name: "single, incomplete frame",
			packet: []byte{
				0x01, // Type
				0xba, // ID
				0x27, // Seq
				// Size header missing
			},
			assert: func(t *testing.T, frames []arnetworkal.Frame, err error) {
				assert.NotNil(t, err)
				assert.Empty(t, frames)
			},
		},
		{
			name: "one good frame, one with a missing byte",
			packet: []byte{
				// Start first frame
				0x01,                   // Type
				0xba,                   // ID
				0x27,                   // Seq
				0x08, 0x00, 0x00, 0x00, // Frame size (little endian)
				0x42, // Data
				// Start second frame
				0x02,                   // Type
				0x0b,                   // ID
				0xc3,                   // Seq
				0x0b, 0x00, 0x00, 0x00, // Frame size (little endian)
				0x12, 0x34, 0x56, // Data (with one byte missing)
			},
			assert: func(t *testing.T, frames []arnetworkal.Frame, err error) {
				assert.NotNil(t, err)
				assert.Empty(t, frames)
			},
		},
		{
			name:   "empty packet",
			packet: []byte{},
			assert: func(t *testing.T, frames []arnetworkal.Frame, err error) {
				assert.Nil(t, err)
				assert.Empty(t, frames)
			},
		},
		{
			name: "single frame",
			packet: []byte{
				0x01,                   // Type
				0xba,                   // ID
				0x27,                   // Seq
				0x08, 0x00, 0x00, 0x00, // Frame size (little endian)
				0x42, // Data
			},
			assert: func(t *testing.T, frames []arnetworkal.Frame, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, len(frames))
			},
		},
		{
			name: "multiple frames",
			packet: []byte{
				// Start first frame
				0x01,                   // Type
				0xba,                   // ID
				0x27,                   // Seq
				0x08, 0x00, 0x00, 0x00, // Frame size (little endian)
				0x42, // Data
				// Start second frame
				0x02,                   // Type
				0x0b,                   // ID
				0xc3,                   // Seq
				0x0b, 0x00, 0x00, 0x00, // Frame size (little endian)
				0x12, 0x34, 0x56, 0x78, // Data
			},
			assert: func(t *testing.T, frames []arnetworkal.Frame, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, len(frames))
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			frames, err := defaultDecodePacket(testCase.packet)
			testCase.assert(t, frames, err)
		})
	}
}
