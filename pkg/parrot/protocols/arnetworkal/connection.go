package arnetworkal

// Connection is an interface implemented by any component capable of
// establishing and maintaining network connections (of some type) over which
// it can both send and receive ARNetworkAL frames.
type Connection interface {
	// Send delivers an ARNetworkAL frame over a network connection of some
	// type.
	Send(frame Frame) error
	// Receive receives 1 or more ARNetworkAL frames over a network connection
	// of some type. Calls to this functions should expect to block until data
	// is received.
	Receive() ([]Frame, error)
	// Close closes any underlying network connections.
	Close()
}
