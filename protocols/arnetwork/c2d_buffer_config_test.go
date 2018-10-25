package arnetwork

import (
	"testing"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/stretchr/testify/require"
)

func TestValidateC2DBufferConfig(t *testing.T) {
	testCases := []struct {
		name       string
		bufCfg     C2DBufferConfig
		assertions func(*testing.T, error)
	}{

		{
			name: "invalid frame type",
			bufCfg: C2DBufferConfig{
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
			bufCfg: C2DBufferConfig{
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
			bufCfg: C2DBufferConfig{
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
