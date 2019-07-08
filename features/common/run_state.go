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
	runID *string
	lock  sync.RWMutex
}

func (r *runState) ID() uint8 {
	return 30
}

func (r *runState) Name() string {
	return "RunState"
}

func (r *runState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"RunIdChanged",
			[]interface{}{
				string(0), // runId,
			},
			r.runIDChanged,
		),
	}
}

// runIDChanged is invoked by the device to provide a unique identifier for the
// current flight.
func (r *runState) runIDChanged(args []interface{}) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.runID = ptr.ToString(args[0].(string))
	log.WithField(
		"runID", *r.runID,
	).Debug("run id changed")
	return nil
}

func (r *runState) RLock() {
	r.lock.RLock()
}

func (r *runState) RUnlock() {
	r.lock.RUnlock()
}

func (r *runState) RunID() (string, bool) {
	if r.runID == nil {
		return "", false
	}
	return *r.runID, true
}
