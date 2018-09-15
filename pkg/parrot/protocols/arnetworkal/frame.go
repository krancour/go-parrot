package arnetworkal

// FrameType is a type for constants used to indicate specific ARNetworkAL
// protocol frame (data) types.
type FrameType uint8

const (
	// frameTypeUninitialized is the zero value for type FrameType and
	// represents an unspecified frame type. This should not be deliberately
	// used!
	frameTypeUninitialized FrameType = 0
	// FrameTypeAck represents a frame that acknowledges receipt of another
	// frame.
	FrameTypeAck FrameType = 1
	// FrameTypeData represents a frame containing data and NOT requiring any
	// acknowledgement of receipt.
	FrameTypeData FrameType = 2
	// FrameTypeLowLatencyData represents a frame containing HIGH PRIORITY data
	// and NOT requiring any acknowledgement of receipt. Such frames are not
	// distinguished from any other frames as they traverse the network, but
	// internally, implementations of the ARNetworkAL protocol prioritize such
	// frames.
	FrameTypeLowLatencyData FrameType = 3
	// FrameTypeDataWithAck represents a frame containing data AND requiring
	// acknowledgement of receipt.
	FrameTypeDataWithAck FrameType = 4
)

// Frame represents a single frame of data to be delivered using the ARNetworkAL
// protocol. Note the data that is actually transferred over the wire is a
// binary encoding of this. Further note that the encoding used is not
// implemented in this package, as it is specific to the Connection
// implementation used for Frame delivery and receipt.
type Frame struct {
	Type FrameType
	ID   uint8
	Seq  uint8
	Data []byte
}
