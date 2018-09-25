package arnetwork

import (
	"time"
)

// OutBufferConfig represents the configuration of a buffer for outbound frames.
// nolint: lll
type OutBufferConfig struct {
	BaseBufferConfig
	AckTimeout time.Duration // Time before considering a frame lost
	MaxRetries int           // Number of retries before considering a frame lost
}
