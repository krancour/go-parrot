package fake

import "github.com/krancour/go-parrot/protocols/arnetworkal"

// FrameSender is a fake implementation of the arnetworkal.FrameSender interface
// used to facilitate unit testing.
type FrameSender struct {
	SendBehavior  func(arnetworkal.Frame) error
	CloseBehavior func()
}

// Send simulates the sending of an arnetworkal.Frame using optional,
// test-defined behavior.
func (f *FrameSender) Send(frame arnetworkal.Frame) error {
	if f.SendBehavior != nil {
		return f.SendBehavior(frame)
	}
	return nil
}

// Close simulates the closing of a connection using optional, test-defined
// behavior.
func (f *FrameSender) Close() {
	if f.CloseBehavior != nil {
		f.CloseBehavior()
	}
}
