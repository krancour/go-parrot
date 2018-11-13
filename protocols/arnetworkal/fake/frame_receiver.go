package fake

import "github.com/krancour/go-parrot/protocols/arnetworkal"

// FrameReceiver is a fake implementation of the arnetworkal.FrameReceiver
// interface used to facilitate unit testing.
type FrameReceiver struct {
	ReceiveBehavior func() ([]arnetworkal.Frame, error)
	CloseBehavior   func()
}

// Receive simulates the receipt of arnetworkal.Frames using optional,
// test-defined behavior.
func (f *FrameReceiver) Receive() ([]arnetworkal.Frame, error) {
	if f.ReceiveBehavior != nil {
		return f.ReceiveBehavior()
	}
	return nil, nil
}

// Close simulates the closing of a connection using optional, test-defined
// behavior.
func (f *FrameReceiver) Close() {
	if f.CloseBehavior != nil {
		f.CloseBehavior()
	}
}
