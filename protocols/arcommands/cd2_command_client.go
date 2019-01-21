package arcommands

import (
	"bytes"
	"encoding/binary"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetwork"
	"github.com/pkg/errors"
)

// BufferType ...
// TODO: Document this
type BufferType uint8

const (
	// BufferTypeAck ...
	// TODO: Document this
	BufferTypeAck BufferType = 1
	// BufferTypeNonAck ...
	// TODO: Document this
	BufferTypeNonAck BufferType = 2
	// BufferTypeHighPrio ...
	// TODO: Document this
	BufferTypeHighPrio BufferType = 3
)

// RetryPolicy ...
// TODO: Document this
type RetryPolicy uint8

const (
	// RetryPolicyPop ...
	// TODO: Document this
	RetryPolicyPop RetryPolicy = 1
	// RetryPolicyRetry ...
	// TODO: Document this
	RetryPolicyRetry RetryPolicy = 2
	// RetryPolicyFlush ...
	// TODO: Document this
	RetryPolicyFlush RetryPolicy = 3
)

// C2DCommandClient ...
// TODO: Document this
type C2DCommandClient interface {
	SendCommand(
		featureID uint8,
		classID uint8,
		commandID uint16,
		bufferType BufferType,
		retryPolicy RetryPolicy,
		args ...interface{},
	) error
}

type c2dCommandClient struct {
	ackCh      chan<- arnetwork.Frame
	nonAckCh   chan<- arnetwork.Frame
	highPrioCh chan<- arnetwork.Frame
}

// NewC2DCommandClient ...
// TODO: Document this
func NewC2DCommandClient(
	ackCh chan<- arnetwork.Frame,
	nonAckCh chan<- arnetwork.Frame,
	highPrioCh chan<- arnetwork.Frame,
) C2DCommandClient {
	return &c2dCommandClient{
		ackCh:      ackCh,
		nonAckCh:   nonAckCh,
		highPrioCh: highPrioCh,
	}
}

// SendC2DCommand ...
// TODO: Document this
func (c *c2dCommandClient) SendCommand(
	featureID uint8,
	classID uint8,
	commandID uint16,
	bufferType BufferType,
	retryPolicy RetryPolicy,
	args ...interface{},
) error {
	encodedData, err := encodeCommand(featureID, classID, commandID, args)
	if err != nil {
		return errors.Wrap(err, "error encoding c2d command")
	}
	frame := arnetwork.Frame{
		Data: encodedData,
		// TODO: Do something with retryPolicy
	}
	switch bufferType {
	case BufferTypeAck:
		c.ackCh <- frame
	case BufferTypeNonAck:
		c.nonAckCh <- frame
	case BufferTypeHighPrio:
		c.highPrioCh <- frame
	default:
		return errors.Errorf("unrecognized buffer type %d", bufferType)
	}
	return nil
}

func encodeCommand(
	featureID uint8,
	classID uint8,
	commandID uint16,
	args ...interface{},
) ([]byte, error) {
	log.Debug("encoding c2d command")
	var buf bytes.Buffer
	buf.WriteByte(featureID)
	buf.WriteByte(classID)
	if err := binary.Write(&buf, binary.LittleEndian, commandID); err != nil {
		return nil, errors.Wrap(err, "error encoding c2d command")
	}
	for _, arg := range args {
		if err := binary.Write(&buf, binary.LittleEndian, arg); err != nil {
			return nil, errors.Wrap(err, "error encoding c2d command")
		}
	}
	return buf.Bytes(), nil
}
