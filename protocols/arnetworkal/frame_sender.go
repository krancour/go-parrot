package arnetworkal

// FrameSender is an interface implemented by any component capable of
// maintaining a network connection (of some type) over which it can send
// ARNetworkAL frames.
type FrameSender interface {
	// Send delivers an ARNetworkAL frame over a network connection of some
	// type.
	Send(Frame) error
	// Close closes the underlying network connection.
	Close()
}
