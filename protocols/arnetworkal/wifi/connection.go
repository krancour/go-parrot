package wifi

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

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
	decodeData func(data []byte) ([]arnetworkal.Frame, error)
}

// NewConnection returns a UDP/IP based implementation of the
// arnetworkal.Connection interface.
func NewConnection() (arnetworkal.Connection, error) {
	// Select an available port
	d2cPort, err := freeport.GetFreePort()
	if err != nil {
		return nil,
			fmt.Errorf("error selecting available client-side port: %s", err)
	}

	// Negotiate the connection. This is how the client informs the device
	// of the UDP port it will listen on. In response, the device informs
	// the client of which UDP port it will listen on.
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
	jsonBytes = append(jsonBytes, 0x00)
	if _, err := negConn.Write(jsonBytes); err != nil {
		return nil,
			fmt.Errorf(
				"error sending connection negotiation request: %s",
				err,
			)
	}
	data, err := bufio.NewReader(negConn).ReadBytes(0x00)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error receiving connection negotiation response: %s",
				err,
			)
	}
	var negRes connectionNegotiationResponse
	if err := json.Unmarshal(data[:len(data)-1], &negRes); err != nil {
		return nil,
			fmt.Errorf(
				"error unmarshaling connection negotiation response: %s",
				err,
			)
	}
	// Any non-zero status is a refused connection.
	if negRes.Status != 0 {
		return nil,
			errors.New(
				"connection negotiation failed; connection refused by device",
			)
	}

	// Establish an outbound connection...
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
		c2dPort:     negRes.C2DPort,
		c2dAddr:     c2dAddr,
		c2dConn:     c2dConn,
		d2cPort:     d2cPort,
		d2cConn:     d2cConn,
		encodeFrame: defaultEncodeFrame,
		decodeData:  defaultDecodeData,
	}, nil
}

func (c *connection) Send(frame arnetworkal.Frame) error {
	if _, err := c.c2dConn.Write(c.encodeFrame(frame)); err != nil {
		return fmt.Errorf("error writing frame to outbound connection: %s", err)
	}
	return nil
}

func (c *connection) Receive() ([]arnetworkal.Frame, error) {
	// TODO: This is probably a huge over-allocation! 64k per frame? Really?
	data := make([]byte, maxUDPDataBytes)
	bytesRead, _, err := c.d2cConn.ReadFromUDP(data)
	if err != nil {
		return nil,
			fmt.Errorf(
				"error receiving frames from inbound connection: %s",
				err,
			)
	}
	return c.decodeData(data[0:bytesRead])
}

func (c *connection) Close() {
	if err := c.c2dConn.Close(); err != nil {
		log.Printf("error closing outbound connection: %s\n", err)
	}
	if err := c.d2cConn.Close(); err != nil {
		log.Printf("error closing inbound connection: %s\n", err)
	}
}
