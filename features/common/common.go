package common

import "github.com/krancour/go-parrot/protocols/arcommands"

// Common ...
// TODO: Document this
type Common interface {
	AllStates() error
}

type common struct {
	c2dCommandClient arcommands.C2DCommandClient
}

func (c *common) AllStates() error {
	return c.c2dCommandClient.SendCommand(
		0, // Common
		4, // Common
		0, // AllStates
		arcommands.BufferTypeAck,
		arcommands.RetryPolicyRetry,
	)
}
