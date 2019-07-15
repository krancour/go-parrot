package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Media streaming status.

// MediaStreamingState ...
// TODO: Document this
type MediaStreamingState interface {
	lock.ReadLockable
}

type mediaStreamingState struct {
	sync.RWMutex
}

func (m *mediaStreamingState) ID() uint8 {
	return 22
}

func (m *mediaStreamingState) Name() string {
	return "MediaStreamingState"
}

func (m *mediaStreamingState) D2CCommands(
	log *log.Entry,
) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"VideoEnableChanged",
			[]interface{}{
				int32(0), // enabled,
			},
			m.videoEnableChanged,
			log,
		),
		arcommands.NewD2CCommand(
			1,
			"VideoStreamModeChanged",
			[]interface{}{
				int32(0), // mode,
			},
			m.videoStreamModeChanged,
			log,
		),
	}
}

// TODO: Implement this
// Title: Video stream state
// Description: Video stream state.
// Support: 0901;090c;090e
// Triggered: by [EnableOrDisableVideoStream](#1-21-0).
// Result:
func (m *mediaStreamingState) videoEnableChanged(
	args []interface{},
	log *log.Entry,
) error {
	// enabled := args[0].(int32)
	//   Current video streaming status.
	//   0: enabled: Video streaming is enabled.
	//   1: disabled: Video streaming is disabled.
	//   2: error: Video streaming failed to start.
	log.Warn("command not implemented")
	return nil
}

// TODO: Implement this
func (m *mediaStreamingState) videoStreamModeChanged(
	args []interface{},
	log *log.Entry,
) error {
	// mode := args[0].(int32)
	//   stream mode
	//   0: low_latency: Minimize latency with average reliability (best for
	//      piloting).
	//   1: high_reliability: Maximize the reliability with an average latency
	//      (best when streaming quality is important but not the latency).
	//   2: high_reliability_low_framerate: Maximize the reliability using a
	//      framerate decimation with an average latency (best when streaming
	//      quality is important but not the latency).
	log.Warn("command not implemented")
	return nil
}
