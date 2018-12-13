package ardrone3

import (
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
)

// Pro features enabled on the Bebop

// PROState ...
// TODO: Document this
type PROState interface{}

type proState struct{}

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

// TODO: Implement this
// Title: Pro features
// Description: Pro features.
// Support:
// Triggered:
// Result:
// WARNING: Deprecated
func (p *proState) features(args []interface{}) error {
	// features := args[0].(uint64)
	//   Bitfield representing enabled features.
	log.Info("ardrone3.features() called")
	return nil
}
