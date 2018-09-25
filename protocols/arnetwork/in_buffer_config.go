package arnetwork

// InBufferConfig represents the configuration of a buffer for inbound frames.
type InBufferConfig struct {
	BaseBufferConfig
	CallBack func(Frame)
}
