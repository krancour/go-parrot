package fake

import "github.com/krancour/go-parrot/protocols/arnetworkal"

// Connection is a fake implementation of the arnetworkal.Connection interface
// used to facilitate unit testing.
type Connection struct {
	SendBehavior    func(arnetworkal.Frame) error
	ReceiveBehavior func() ([]arnetworkal.Frame, error)
	CloseBehavior   func()
}

// Send simulates the sending of an arnetworkal.Frame using optional,
// test-defined behavior.
func (c *Connection) Send(frame arnetworkal.Frame) error {
	if c.SendBehavior != nil {
		return c.SendBehavior(frame)
	}
	return nil
}

// Receive simulates the receipt of arnetworkal.Frames using optional,
// test-defined behavior.
func (c *Connection) Receive() ([]arnetworkal.Frame, error) {
	if c.ReceiveBehavior != nil {
		return c.ReceiveBehavior()
	}
	return nil, nil
}

// Close simulates the closing of a connection using optional, test-defined
// behavior.
func (c *Connection) Close() {
	if c.CloseBehavior != nil {
		c.CloseBehavior()
	}
}
