package arnetwork

import (
	"testing"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/stretchr/testify/require"
)

func TestValidateD2CBufferConfig(t *testing.T) {
	testCases := []struct {
		name       string
		bufCfg     D2CBufferConfig
		assertions func(*testing.T, error)
	}{

		{
			name: "invalid frame type",
			bufCfg: D2CBufferConfig{
				FrameType: 0, // This frame type is undefined
				Size:      1,
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid frame type")
			},
		},

		{
			name: "invalid size",
			bufCfg: D2CBufferConfig{
				FrameType: arnetworkal.FrameTypeData,
				Size:      0,
			},
			assertions: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "invalid size")
			},
		},

		{
			name: "valid config",
			bufCfg: D2CBufferConfig{
				FrameType: arnetworkal.FrameTypeData,
				Size:      1,
			},
			assertions: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.bufCfg.validate()
			testCase.assertions(t, err)
		})
	}
}
