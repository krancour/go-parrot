package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Sounds related events

// SoundState ...
// TODO: Document this
type SoundState interface {
	lock.ReadLockable
}

type soundState struct {
	sync.RWMutex
}

func newSoundState() *soundState {
	return &soundState{}
}

func (s *soundState) ID() uint8 {
	return 36
}

func (s *soundState) Name() string {
	return "SoundState"
}

func (s *soundState) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"AlertSound",
			[]interface{}{
				int32(0), // state,
			},
			s.alertSound,
			log,
		),
	}
}

// TODO: Implement this
// Title: Alert sound state
// Description: Alert sound state.
// Support: none
// Triggered: by [StartAlertSound](#1-35-0) or [StopAlertSound](#1-35-1) or \
//   when the drone starts or stops to play an alert sound by itself.
// Result:
func (s *soundState) alertSound(args []interface{}, log *log.Entry) error {
	// state := args[0].(int32)
	//   State of the alert sound
	//   0: stopped: Alert sound is not playing
	//   1: playing: Alert sound is playing
	log.Warn("command not implemented")
	return nil
}
