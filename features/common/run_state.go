package common

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Commands sent by the drone to inform about the run or flight state

// RunState ...
// TODO: Document this
type RunState interface{}

type runState struct{}

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
			r.runIdChanged,
		),
	}
}

// TODO: Implement this
// Title: Current run id
// Description: Current run id.\n A run id is uniquely identifying a run or a
//   flight.\n For each run is generated on the drone a file which can be used
//   by Academy to sum up the run.\n Also, each medias taken during a run has a
//   filename containing the run id.
// Support: 0901:3.0.1;090c;090e
// Triggered: when the drone generates a new run id (generally right after a
//   take off).
// Result:
func (r *runState) runIdChanged(args []interface{}) error {
	// runId := args[0].(string)
	//   Id of the run
	log.Info("common.runIdChanged() called")
	return nil
}
