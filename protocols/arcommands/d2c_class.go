package arcommands

import log "github.com/Sirupsen/logrus"

// D2CClass ...
// TODO: Document this
type D2CClass interface {
	ClassID() uint8
	ClassName() string
	D2CCommands(log *log.Entry) []D2CCommand
}
