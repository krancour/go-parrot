package wifi

import (
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/pkg/errors"
)

type frameSender struct {
	conn *net.UDPConn
	// This function is overridable by unit tests
	encodeFrame func(frame arnetworkal.Frame) ([]byte, error)
}

func (f *frameSender) Send(frame arnetworkal.Frame) error {
	log := log.WithField(
		"uuid", frame.UUID,
	).WithField(
		"buffer", frame.ID,
	).WithField(
		"type", frame.Type,
	).WithField(
		"seq", frame.Seq,
	)
	frameBytes, err := f.encodeFrame(frame)
	if err != nil {
		return errors.Wrap(err, "error encoding arnetworkal frame")
	}
	log.Debug("sending arnetworkal frame")
	if _, err := f.conn.Write(frameBytes); err != nil {
		return errors.Wrap(
			err,
			"error writing datagram to c2d connection",
		)
	}
	log.Debug("sent arnetworkal frame")
	return nil
}

func (f *frameSender) Close() {
	if f.conn != nil {
		log.Debug("closing c2d connection")
		if err := f.conn.Close(); err != nil {
			log.Errorf("error closing c2d connection: %s", err)
		}
		log.Debug("closed c2d connection")
	}
}
