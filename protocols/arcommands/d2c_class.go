package arcommands

import log "github.com/Sirupsen/logrus"

// D2CClass ...
// TODO: Document this
type D2CClass interface {
	ID() uint8
	Name() string
	D2CCommands(log *log.Entry) []D2CCommand
}
