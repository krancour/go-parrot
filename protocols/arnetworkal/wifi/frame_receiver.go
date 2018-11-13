package wifi

import (
	"net"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/pkg/errors"
)

type frameReceiver struct {
	conn *net.UDPConn
	// This function is overridable by unit tests
	decodeDatagram     func(data []byte) ([]arnetworkal.Frame, error)
	datagramBuffer     []byte
	datagramBufferLock sync.Mutex
}

func (f *frameReceiver) Receive() ([]arnetworkal.Frame, error) {
	f.datagramBufferLock.Lock()
	defer f.datagramBufferLock.Unlock()
	log.Debug("reading / waiting for datagram from d2c connection")
	bytesRead, _, err := f.conn.ReadFromUDP(f.datagramBuffer)
	if err != nil {
		return nil,
			errors.Wrap(err, "error receiving datagram from d2c connection")
	}
	log.WithField(
		"bytesRead", bytesRead,
	).Debug("got datagram from d2c connection")
	return f.decodeDatagram(f.datagramBuffer[0:bytesRead])
}

func (f *frameReceiver) Close() {
	if f.conn != nil {
		log.Debug("closing d2c connection")
		if err := f.conn.Close(); err != nil {
			log.Errorf("error closing d2c connection: %s", err)
		}
		log.Debug("closed d2c connection")
	}
}
