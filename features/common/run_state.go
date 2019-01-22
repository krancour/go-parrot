package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Commands sent by the drone to inform about the run or flight state

// RunState ...
// TODO: Document this
type RunState interface {
	// RLock blocks until a read lock is obtained. This permits callers to procede
	// with querying any or all attributes of the run state without worry
	// that some attributes will be overwritten as others are read. i.e. It
	// permits the possibility of taking an atomic snapshop of run state.
	// Note that use of this function is not obligatory for applications that do
	// not require such guarantees. Callers MUST call RUnlock() or else run
	// state will never resume updating.
	RLock()
	// RUnlock releases a read lock on the run state. See RLock().
	RUnlock()
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
