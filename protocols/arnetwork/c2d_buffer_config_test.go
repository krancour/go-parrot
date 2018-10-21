package arnetwork

import (
	"testing"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/stretchr/testify/require"
)

func TestInvalidC2DBufferConfig(t *testing.T) {
	bufCfg := C2DBufferConfig{
		FrameType: 5, // This is an invalid frame type
	}
	err := bufCfg.validate()
	require.Error(t, err)
}

func TestValidC2DBufferConfig(t *testing.T) {
	bufCfg := C2DBufferConfig{
		FrameType: arnetworkal.FrameTypeData,
	}
	err := bufCfg.validate()
	require.NoError(t, err)
}
