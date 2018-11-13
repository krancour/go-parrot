package arcommands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeArgs(t *testing.T) {
	data := []byte{
		1,    // FeatureID
		4,    // ClassID
		9, 0, // CommandID (little endian)
		1,    // uint8
		2,    // int8
		3, 0, // uint16 (little endian)
		4, 0, // int16 (little endian)
		5, 0, 0, 0, // uint32 (little endian)
		6, 0, 0, 0, // int32 (little endian)
		7, 0, 0, 0, 0, 0, 0, 0, // uint64 (little endian)
		8, 0, 0, 0, 0, 0, 0, 0, // int164 (little endian)
		102, 111, 111, 0, // null terminated string
		154, 153, 17, 65, // float32 (little endian)
		102, 102, 102, 102, 102, 102, 36, 64, // float64 (little endian)
	}
	args := []interface{}{
		uint8(0),
		int8(0),
		uint16(0),
		int16(0),
		uint32(0),
		int32(0),
		uint64(0),
		int64(0),
		"",
		float32(0),
		float64(0),
	}
	expected := []interface{}{
		uint8(1),
		int8(2),
		uint16(3),
		int16(4),
		uint32(5),
		int32(6),
		uint64(7),
		int64(8),
		"foo",
		float32(9.1),
		float64(10.2),
	}
	err := decodeArgs(data, args)
	require.NoError(t, err)
	assert.Equal(t, expected, args)
}
