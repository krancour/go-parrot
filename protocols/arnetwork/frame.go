package arnetwork

// Frame represents a frame to be sent or received over the arnetwork
// protocol. At this level of abstraction, no details of the underlying
// network protocols (arnetworkal, UDP/IP, BLE, etc.) bleed through.
type Frame struct {
	Data []byte
	seq  uint8
}
