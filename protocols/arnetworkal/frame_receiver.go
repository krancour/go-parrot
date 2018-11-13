package arnetworkal

// FrameReceiver is an interface implemented by any component capable of
// maintaining a network connection (of some type) over which it can receive
// ARNetworkAL frames.
type FrameReceiver interface {
	// Receive receives 1 or more ARNetworkAL frames over a network connection
	// of some type. Calls to this functions should expect to block until data
	// is received.
	Receive() ([]Frame, error)
	// Close closes the underlying network connection.
	Close()
}
