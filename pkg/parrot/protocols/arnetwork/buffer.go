package arnetwork

// TODO: Document this
type Buffer interface {
}

// TODO: Do I want to change some of these attribute names or types?
type buffer struct {
	id              int
	dataType        arnetworkal.FrameType
	numberOfRetry   int
	numberOfCell    int32
	dataCopyMaxSize int32
	isOverwriting   int // TODO: Should this be a bool?
}

// TODO: Document this
func NewBuffer() Buffer {
	// TODO: Implement the rest of this
	return &buffer{}
}
