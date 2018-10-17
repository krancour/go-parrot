package wifi

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/phayes/freeport"
)

const (
	// maxUDPDataBytes represents the practical maximum numbers of data bytes in
	// a UDP packet.
	maxUDPDataBytes = 65507
)

var (
	// These are vars instead of a consts so that they can be overridden by unit
	// tests. They're not exported, so there is no danger of anyone else
	// tampering with these.
	deviceIP      = "192.168.42.1"
	discoveryPort = 44444
)

type connectionNegotiationRequest struct {
	D2CPort        int    `json:"d2c_port"`
	ControllerType string `json:"controller_type"`
	ControllerName string `json:"controller_name"`
}

type connectionNegotiationResponse struct {
	Status  int `json:"status"`
	C2DPort int `json:"c2d_port"`
}

type connection struct {
	c2dPort int
	c2dAddr *net.UDPAddr
	c2dConn *net.UDPConn
	d2cPort int
	d2cConn *net.UDPConn
	// This function is overridable by unit tests
	encodeFrame func(frame arnetworkal.Frame) []byte
	// This function is overridable by unit tests
	decodePacket func(data []byte) ([]arnetworkal.Frame, error)
}

// NewConnection returns a UDP/IP based implementation of the
// arnetworkal.Connection interface.
func NewConnection() (arnetworkal.Connection, error) {
	log.Debug("starting new connection process")

	// Select an available port
	log.Debug("selecting available port for d2c communication")
	d2cPort, err := freeport.GetFreePort()
	if err != nil {
		return nil,
			fmt.Errorf("error selecting available client-side port: %s", err)
	}
	log.WithField(
		"port", d2cPort,
	).Debug("selected port for d2c communication")

	// Negotiate the connection. This is how the client informs the device
	// of the UDP port it will listen on. In response, the device informs
	// the client of which UDP port it will listen on.
	// TODO: Should this be moved into its own protocol packages?
	log.WithField(
		"deviceIP", deviceIP,
	).WithField(
		"discoveryPort", discoveryPort,
	).Debug("negotiating connection")
	negAddr, err := net.ResolveTCPAddr(
		"tcp",
		fmt.Sprintf("%s:%d", deviceIP, discoveryPort),
	)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error resolving address for connection negotiation: %s",
				err,
			)
	}
	negConn, err := net.DialTCP("tcp", nil, negAddr)
	if err != nil {
		return nil, fmt.Errorf("error negotiating connection: %s", err)
	}
	defer negConn.Close()
	log.Debug("marshaling connection negotiation request")
	jsonBytes, err := json.Marshal(
		connectionNegotiationRequest{
			D2CPort:        d2cPort,
			ControllerType: "computer",
			ControllerName: "go-parrot",
		},
	)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error marshaling connection negotiation request: %s",
				err,
			)
	}
	log.Debug("marshaled connection negotiation request")
	jsonBytes = append(jsonBytes, 0x00)
	log.Debug("sending connection negotiation request")
	if _, err := negConn.Write(jsonBytes); err != nil {
		return nil,
			fmt.Errorf(
				"error sending connection negotiation request: %s",
				err,
			)
	}
	log.Debug("sent connection negotiation request")
	log.Debug("waiting for connection negotiation response")
	// Note: Since TCP is a stream-based protocol, we either need to know
	// content length of the response in advance OR a delimiter we can look for.
	// Since we know the remote end of the connection is implemented in C, it's
	// pretty reasonable to use the null character 0x00 (which is used to
	// terminate all strings in C) as a delimiter.
	data, err := bufio.NewReader(negConn).ReadBytes(0x00)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error receiving connection negotiation response: %s",
				err,
			)
	}
	log.Debug("got connection negotiation response")
	var negRes connectionNegotiationResponse
	log.Debug("unmarshaling connection negotiation response")
	if err := json.Unmarshal(data[:len(data)-1], &negRes); err != nil {
		return nil,
			fmt.Errorf(
				"error unmarshaling connection negotiation response: %s",
				err,
			)
	}
	log.Debug("unmarshaled connection negotiation response")
	// Any non-zero status is a refused connection.
	if negRes.Status != 0 {
		return nil,
			errors.New(
				"connection negotiation failed; connection refused by device",
			)
	}

	log.WithField(
		"deviceIP", deviceIP,
	).WithField(
		"c2dPort", negRes.C2DPort,
	).WithField(
		"d2cPort", d2cPort,
	).Debug("connection negotiation complete")

	// Establish an outbound connection...
	log.WithField(
		"deviceIP", deviceIP,
	).WithField(
		"c2dPort", negRes.C2DPort,
	).Debug("establishing c2d connection")
	c2dAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf("%s:%d", deviceIP, negRes.C2DPort),
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
	log.WithField(
		"deviceIP", deviceIP,
	).WithField(
		"c2dPort", negRes.C2DPort,
	).Debug("established c2d connection")

	// Establish an inbound connection...
	log.WithField(
		"d2cPort", d2cPort,
	).Debug("establishing d2c connection")
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
	log.WithField(
		"d2cPort", d2cPort,
	).Debug("established d2c connection")

	log.WithField(
		"deviceIP", deviceIP,
	).WithField(
		"c2dPort", negRes.C2DPort,
	).WithField(
		"d2cPort", d2cPort,
	).Debug("c2d and d2c connections ready for use")
	return &connection{
		c2dPort:      negRes.C2DPort,
		c2dAddr:      c2dAddr,
		c2dConn:      c2dConn,
		d2cPort:      d2cPort,
		d2cConn:      d2cConn,
		encodeFrame:  defaultEncodeFrame,
		decodePacket: defaultDecodePacket,
	}, nil
}

func (c *connection) Send(frame arnetworkal.Frame) error {
	log := log.WithField(
		"buffer", frame.ID,
	).WithField(
		"type", frame.Type,
	).WithField(
		"seq", frame.Seq,
	)
	log.Debug("sending arnetworkal frame")
	if _, err := c.c2dConn.Write(c.encodeFrame(frame)); err != nil {
		return fmt.Errorf("error writing frame to outbound connection: %s", err)
	}
	log.Debug("sent arnetworkal frame")
	return nil
}

func (c *connection) Receive() ([]arnetworkal.Frame, error) {
	log.Debug("reading / waiting for datagram from d2c connection")
	// TODO: We should be able to use this slice over and over again for
	// efficiency
	packet := make([]byte, maxUDPDataBytes)
	bytesRead, _, err := c.d2cConn.ReadFromUDP(packet)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error receiving frames from inbound connection: %s",
				err,
			)
	}
	log.WithField(
		"bytesRead", bytesRead,
	).Debug("got datagram from d2c connection")
	return c.decodePacket(packet[0:bytesRead])
}

func (c *connection) Close() {
	log.Debug("closing c2d connection")
	if err := c.c2dConn.Close(); err != nil {
		log.Printf("error closing outbound connection: %s\n", err)
	}
	log.Debug("closed c2d connection")
	log.Debug("closing d2c connection")
	if err := c.d2cConn.Close(); err != nil {
		log.Printf("error closing inbound connection: %s\n", err)
	}
	log.Debug("closed d2c connection")
}
