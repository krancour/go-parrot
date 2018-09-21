package arnetworkal

// Connection is an interface implemented by any component capable of
// establishing and maintaining network connections (of some type) over which
// it can both send and receive ARNetworkAL frames.
type Connection interface {
	// Send delivers an ARNetworkAL frame over a network connection of some
	// type.
	Send(Frame) error
	// Receive attempts to receive one ARNetworkAL frame from an underlying
	// network connection of some type. It also returns a boolean indicating
	// success or failure since it is possible to for this function to return
	// without error if the underlying connection has been closed. It also returns
	// any error that is encountered.
	Receive() (Frame, bool, error)
	// Close closes any underlying network connections.
	Close() error
}
