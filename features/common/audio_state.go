package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Audio-related state updates.

// AudioState ...
// TODO: Document this
type AudioState interface{}

type audioState struct{}

func (a *audioState) ID() uint8 {
	return 21
}

func (a *audioState) Name() string {
	return "AudioState"
}

func (a *audioState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"AudioStreamingRunning",
			[]interface{}{
				uint8(0), // running,
			},
			a.audioStreamingRunning,
		),
	}
}

// TODO: Implement this
// Title: Audio stream direction
// Description: Audio stream direction.
// Support: 0905;0906
// Triggered: by [SetAudioStreamDirection](#0-20-0).
// Result:
func (a *audioState) audioStreamingRunning(args []interface{}) error {
	// running := args[0].(uint8)
	//   Bit field for TX and RX running bit 0 is 1 if Drone TX is running bit 1
	//   is 1 if Drone RX is running
	log.Info("common.audioStreamingRunning() called")
	return nil
}
