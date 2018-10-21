package arnetwork

import (
	"testing"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/stretchr/testify/require"
)

func TestInvalidD2cBufferConfig(t *testing.T) {
	bufCfg := D2CBufferConfig{
		FrameType: 5, // This is an invalid frame type
	}
	err := bufCfg.validate()
	require.Error(t, err)
}

func TestValidD2CBufferConfig(t *testing.T) {
	bufCfg := D2CBufferConfig{
		FrameType: arnetworkal.FrameTypeData,
	}
	err := bufCfg.validate()
	require.NoError(t, err)
}
