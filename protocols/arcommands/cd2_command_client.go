package arcommands

// C2DCommandClient ...
// TODO: Document this
type C2DCommandClient interface {
	SendCommand(
		featureID uint8,
		classID uint8,
		commandID uint16,
		args ...interface{},
	) error
}

type c2dCommandClient struct{}

// NewC2DCommandClient ...
// TODO: Document this
// TODO: What should we pass to this?
func NewC2DCommandClient() C2DCommandClient {
	return &c2dCommandClient{}
}

// SendC2DCommand ...
// TODO: Document this
// TODO: Implement this
func (c *c2dCommandClient) SendCommand(
	featureID uint8,
	classID uint8,
	commandID uint16,
	args ...interface{},
) error {
	return nil
}
