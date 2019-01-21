package common

import "github.com/krancour/go-parrot/protocols/arcommands"

// Settings ...
// TODO: Document this
type Settings interface {
	AllSettings() error
}

type settings struct {
	c2dCommandClient arcommands.C2DCommandClient
}

func (s *settings) AllSettings() error {
	return s.c2dCommandClient.SendCommand(
		0, // Common
		2, // Settings
		0, // AllSettings
		arcommands.BufferTypeAck,
		arcommands.RetryPolicyRetry,
	)
}
