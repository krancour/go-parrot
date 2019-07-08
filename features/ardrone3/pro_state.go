package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Pro features enabled on the Bebop

// PROState ...
// TODO: Document this
type PROState interface {
	lock.ReadLockable
}

type proState struct {
	lock sync.RWMutex
}

func (p *proState) ID() uint8 {
	return 32
}

func (p *proState) Name() string {
	return "PROState"
}

func (p *proState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"Features",
			[]interface{}{
				uint64(0), // features,
			},
			p.features,
		),
	}
}

// features is deprecated, but since we can still see this command being
// invoked, we'll implement the command to avoid a warning, but the
// implementation will remain a no-op unless / until such time that it becomes
// clear that older versions of the firmware might require us to support it.
func (p *proState) features(args []interface{}) error {
	// features := args[0].(uint64)
	//   Bitfield representing enabled features.
	log.Debug("features changed-- this is a no-op")
	return nil
}

func (p *proState) RLock() {
	p.lock.RLock()
}

func (p *proState) RUnlock() {
	p.lock.RUnlock()
}
