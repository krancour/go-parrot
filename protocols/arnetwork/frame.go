package arnetwork

// Frame ...
// TODO: Document this
// TODO: Not yet positive that this has the right attributes
type Frame struct {
	ID   uint8
	Seq  uint8
	Data []byte
}
