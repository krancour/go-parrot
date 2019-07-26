package common

import (
	"sync"

	"github.com/krancour/go-parrot/lock"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Commands sent by the drone to inform about the run or flight state

// RunState ...
// TODO: Document this
type RunState interface {
	lock.ReadLockable
	// RunID returns a unique identifier for the current flight. A boolean value
	// is also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	RunID() (string, bool)
}

type runState struct {
	sync.RWMutex
	runID *string
}

func newRunState() *runState {
	return &runState{}
}

func (r *runState) ClassID() uint8 {
	return 30
}

func (r *runState) ClassName() string {
	return "RunState"
}

func (r *runState) D2CCommands(log *log.Entry) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"RunIdChanged",
			[]interface{}{
				string(0), // runId,
			},
			r.runIDChanged,
			log,
		),
	}
}

// runIDChanged is invoked by the device to provide a unique identifier for the
// current flight.
func (r *runState) runIDChanged(args []interface{}, log *log.Entry) error {
	r.Lock()
	defer r.Unlock()
	r.runID = ptr.ToString(args[0].(string))
	log.WithField(
		"runID", *r.runID,
	).Debug("run id changed")
	return nil
}

func (r *runState) RunID() (string, bool) {
	if r.runID == nil {
		return "", false
	}
	return *r.runID, true
}
