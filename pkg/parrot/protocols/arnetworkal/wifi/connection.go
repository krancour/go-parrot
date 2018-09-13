package wifi

import (
	"fmt"
	"log"
	"net"

	"github.com/krancour/drone-examples/pkg/parrot/protocols/arnetworkal"
)

// maxUDPDataBytes represents the practical maximum numbers of data bytes in a
// UDP packet.
const maxUDPDataBytes = 65507

var (
	// These are vars instead of a consts so that they can be overridden by unit
	// tests. They're not exported, so there is no danger of anyone else
	// tampering with these.
	deviceIP = "192.168.42.1"
	c2dPort  = 54321
	d2cPort  = 43210
)

type connection struct {
	c2dConn *net.UDPConn
	d2cConn *net.UDPConn
	// This function is overridable by unit tests
	encodeFrame func(frame arnetworkal.Frame) []byte
	// This function is overridable by unit tests
	decodePacket func(data []byte) ([]arnetworkal.Frame, error)
}

// NewConnection returns a UDP/IP based implementation of the
// arnetworkal.Connection interface.
func NewConnection() (arnetworkal.Connection, error) {
	// Establish an outbound connection...
	c2dAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf("%s:%d", deviceIP, c2dPort),
	)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error resolving address for outbound connection: %s",
				err,
			)
	}
	c2dConn, err := net.DialUDP("udp", nil, c2dAddr)
	if err != nil {
		return nil,
			fmt.Errorf("error establishing outbound connection: %s", err)
	}

	// Establish an inbound connection...
	d2cAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf(":%d", d2cPort),
	)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error resolving address for inbound connection: %s",
				err,
			)
	}
	d2cConn, err := net.ListenUDP("udp", d2cAddr)
	if err != nil {
		return nil, fmt.Errorf("error establishing inbound connection: %s", err)
	}

	return &connection{
		c2dConn:      c2dConn,
		d2cConn:      d2cConn,
		encodeFrame:  defaultEncodeFrame,
		decodePacket: defaultDecodePacket,
	}, nil
}

func (c *connection) Send(frame arnetworkal.Frame) error {
	if _, err := c.c2dConn.Write(c.encodeFrame(frame)); err != nil {
		return fmt.Errorf("error writing frame to outbound connection: %s", err)
	}
	return nil
}

func (c *connection) Receive() ([]arnetworkal.Frame, error) {
	packet := make([]byte, maxUDPDataBytes)
	bytesRead, _, err := c.d2cConn.ReadFromUDP(packet)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error receiving frames from inbound connection: %s",
				err,
			)
	}
	return c.decodePacket(packet[0:bytesRead])
}

func (c *connection) Close() {
	if err := c.c2dConn.Close(); err != nil {
		log.Printf("error closing outbound connection: %s\n", err)
	}
	if err := c.d2cConn.Close(); err != nil {
		log.Printf("error closing inbound connection: %s\n", err)
	}
}
