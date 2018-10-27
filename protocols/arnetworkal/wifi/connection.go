package wifi

import (
	"bufio"
	"encoding/json"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/phayes/freeport"
	"github.com/pkg/errors"
)

const (
	// maxUDPDataBytes represents the practical maximum numbers of data bytes in
	// a UDP datagram.
	maxUDPDataBytes = 65507
	discoveryPort   = 44444
)

var (
	deviceIP = net.ParseIP("192.168.42.1")
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
	c2dConn *net.UDPConn
	d2cConn *net.UDPConn
	// This function is overridable by unit tests
	encodeFrame func(frame arnetworkal.Frame) ([]byte, error)
	// This function is overridable by unit tests
	decodeDatagram func(data []byte) ([]arnetworkal.Frame, error)
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
			errors.Wrap(err, "error selecting available client-side port")
	}
	log.WithField(
		"port", d2cPort,
	).Debug("selected port for d2c communication")

	// Negotiate the connection
	c2dPort, err := negotiateConnection(deviceIP, discoveryPort, d2cPort)
	if err != nil {
		return nil, errors.Wrap(err, "connection negotiation failed")
	}

	// Establish the c2d connection...
	c2dConn, err := establishC2DConnection(deviceIP, c2dPort)
	if err != nil {
		return nil, errors.Wrap(err, "error establishing c2d connection")
	}

	// Establish the d2c connection...
	d2cConn, err := establishD2CConnection(d2cPort)
	if err != nil {
		return nil, errors.Wrap(err, "error establishing d2c connection")
	}

	log.WithField(
		"deviceIP", deviceIP,
	).WithField(
		"c2dPort", c2dPort,
	).WithField(
		"d2cPort", d2cPort,
	).Debug("c2d and d2c connections ready for use")
	return &connection{
		c2dConn:        c2dConn,
		d2cConn:        d2cConn,
		encodeFrame:    defaultEncodeFrame,
		decodeDatagram: defaultDecodeDatagram,
	}, nil
}

// negotiateConnection negotiates the connection. This is how the client informs
// the device of the UDP port it will listen on. In response, the device informs
// the client of which UDP port it will listen on.
// TODO: Should this be moved into its own protocol packages?
var defaultNegotiateConnection = func(
	deviceIP net.IP,
	discoveryPort,
	d2cPort int,
) (int, error) {
	log.WithField(
		"deviceIP", deviceIP,
	).WithField(
		"discoveryPort", discoveryPort,
	).Debug("negotiating connection")

	conn, err := net.DialTCP(
		"tcp",
		nil,
		&net.TCPAddr{
			IP:   deviceIP,
			Port: discoveryPort,
		},
	)
	if err != nil {
		return 0, err
	}
	defer conn.Close() // nolint: errcheck

	log.Debug("marshaling connection negotiation request")
	jsonBytes, err := json.Marshal(
		connectionNegotiationRequest{
			D2CPort:        d2cPort,
			ControllerType: "computer",
			ControllerName: "go-parrot",
		},
	)
	if err != nil {
		return 0,
			errors.Wrap(err, "error marshaling connection negotiation request")
	}
	log.Debug("marshaled connection negotiation request")
	// Use a null character to terminate the request
	jsonBytes = append(jsonBytes, 0x00)

	log.Debug("sending connection negotiation request")
	if _, err = conn.Write(jsonBytes); err != nil {
		return 0,
			errors.Wrap(err, "error sending connection negotiation request")
	}
	log.Debug("sent connection negotiation request")

	log.Debug("waiting for connection negotiation response")
	// Note: Since TCP is a stream-based protocol, we either need to know
	// content length of the response in advance OR a delimiter we can look for.
	// Since we know the remote end of the connection is implemented in C, it's
	// pretty reasonable to use the null character 0x00 (which is used to
	// terminate all strings in C) as a delimiter.
	data, err := bufio.NewReader(conn).ReadBytes(0x00)
	if err != nil {
		return 0,
			errors.Wrap(err, "error receiving connection negotiation response")
	}
	log.Debug("got connection negotiation response")

	var res connectionNegotiationResponse
	log.Debug("unmarshaling connection negotiation response")
	if err := json.Unmarshal(data[:len(data)-1], &res); err != nil {
		return 0,
			errors.Wrap(
				err,
				"error unmarshaling connection negotiation response",
			)
	}
	log.Debug("unmarshaled connection negotiation response")
	// Any non-zero status is a refused connection.
	if res.Status != 0 {
		return 0, errors.New("connection refused by device")
	}

	log.WithField(
		"deviceIP", deviceIP,
	).WithField(
		"c2dPort", res.C2DPort,
	).WithField(
		"d2cPort", d2cPort,
	).Debug("connection negotiation complete")

	return res.C2DPort, nil
}

var negotiateConnection = defaultNegotiateConnection

// establishC2DConnection establishes the client to device UDP connection.
var defaultEstablishC2DConnection = func(
	deviceIP net.IP,
	c2dPort int,
) (*net.UDPConn, error) {
	log.WithField(
		"deviceIP", deviceIP,
	).WithField(
		"c2dPort", c2dPort,
	).Debug("establishing c2d connection")

	conn, err := net.DialUDP(
		"udp",
		nil,
		&net.UDPAddr{
			IP:   deviceIP,
			Port: c2dPort,
		},
	)
	if err != nil {
		return nil, err
	}

	log.WithField(
		"deviceIP", deviceIP,
	).WithField(
		"c2dPort", c2dPort,
	).Debug("established c2d connection")

	return conn, nil
}

var establishC2DConnection = defaultEstablishC2DConnection

// establishD2CConnection establishes the device to client UDP connection.
var defaultEstablishD2CConnection = func(d2cPort int) (*net.UDPConn, error) {
	log.WithField(
		"d2cPort", d2cPort,
	).Debug("establishing d2c connection")

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: d2cPort})
	if err != nil {
		return nil, err
	}

	log.WithField(
		"d2cPort", d2cPort,
	).Debug("established d2c connection")

	return conn, nil
}

func (c *connection) Send(frame arnetworkal.Frame) error {
	log := log.WithField(
		"uuid", frame.UUID,
	).WithField(
		"buffer", frame.ID,
	).WithField(
		"type", frame.Type,
	).WithField(
		"seq", frame.Seq,
	)
	frameBytes, err := c.encodeFrame(frame)
	if err != nil {
		return errors.Wrap(err, "error encoding arnetworkal frame")
	}
	log.Debug("sending arnetworkal frame")
	if _, err := c.c2dConn.Write(frameBytes); err != nil {
		return errors.Wrap(
			err,
			"error writing datagram to c2d connection",
		)
	}
	log.Debug("sent arnetworkal frame")
	return nil
}

var establishD2CConnection = defaultEstablishD2CConnection

var datagram = make([]byte, maxUDPDataBytes)

func (c *connection) Receive() ([]arnetworkal.Frame, error) {
	log.Debug("reading / waiting for datagram from d2c connection")
	bytesRead, _, err := c.d2cConn.ReadFromUDP(datagram)
	if err != nil {
		return nil,
			errors.Wrap(err, "error receiving datagram from d2c connection")
	}
	log.WithField(
		"bytesRead", bytesRead,
	).Debug("got datagram from d2c connection")
	return c.decodeDatagram(datagram[0:bytesRead])
}

func (c *connection) Close() {
	c.closeC2DConnection()
	c.closeD2CConnection()
}

func (c *connection) closeC2DConnection() {
	if c.c2dConn != nil {
		log.Debug("closing c2d connection")
		if err := c.c2dConn.Close(); err != nil {
			log.Errorf("error closing c2d connection: %s", err)
		}
		log.Debug("closed c2d connection")
	}
}

func (c *connection) closeD2CConnection() {
	if c.d2cConn != nil {
		log.Debug("closing d2c connection")
		if err := c.d2cConn.Close(); err != nil {
			log.Errorf("error closing d2c connection: %s", err)
		}
		log.Debug("closed d2c connection")
	}
}
