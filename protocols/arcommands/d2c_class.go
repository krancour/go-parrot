package arcommands

// D2CClass ...
// TODO: Document this
type D2CClass interface {
	ID() uint8
	Name() string
	D2CCommands() []D2CCommand
}
