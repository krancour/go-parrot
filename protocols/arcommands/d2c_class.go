package arcommands

// D2CClass ...
// TODO: Document this
type D2CClass interface {
	d2cCommands() map[uint16]D2CCommand
}

type d2cClass struct {
	name    string
	d2cCmds map[uint16]D2CCommand
}

// NewD2CClass ...
// TODO: Document this
func NewD2CClass(name string, d2cCmds map[uint16]D2CCommand) D2CClass {
	return &d2cClass{
		name:    name,
		d2cCmds: d2cCmds,
	}
}

func (d *d2cClass) d2cCommands() map[uint16]D2CCommand {
	return d.d2cCmds
}
