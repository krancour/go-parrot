package wifi

import (
	"bufio"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/phayes/freeport"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConnection(t *testing.T) {
	testCases := []struct {
		name                        string
		negotiateConnectionBehavior func(
			deviceIP net.IP,
			discoveryPort,
			d2cPort int,
		) (int, error)
		establishC2DConnectionBehavior func(
			deviceIP net.IP,
			c2dPort int,
		) (*net.UDPConn, error)
		establishD2CConnectionBehavior func(d2cPort int) (*net.UDPConn, error)
		assertions                     func(*testing.T, arnetworkal.Connection, error)
	}{

		{
			name: "connection negotiation fails",
			negotiateConnectionBehavior: func(net.IP, int, int) (int, error) {
				return 0, errors.New("foo")
			},
			assertions: func(t *testing.T, _ arnetworkal.Connection, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "error negotiating connection")
				assert.Contains(t, err.Error(), "foo")
			},
		},

		{
			name: "establish c2d connection fails",
			negotiateConnectionBehavior: func(net.IP, int, int) (int, error) {
				return 12345, nil
			},
			establishC2DConnectionBehavior: func(net.IP, int) (*net.UDPConn, error) {
				return nil, errors.New("bar")
			},
			assertions: func(t *testing.T, _ arnetworkal.Connection, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "error establishing c2d connection")
				assert.Contains(t, err.Error(), "bar")
			},
		},

		{
			name: "establish d2c connection fails",
			negotiateConnectionBehavior: func(net.IP, int, int) (int, error) {
				return 12345, nil
			},
			establishC2DConnectionBehavior: func(net.IP, int) (*net.UDPConn, error) {
				return nil, nil
			},
			establishD2CConnectionBehavior: func(int) (*net.UDPConn, error) {
				return nil, errors.New("bat")
			},
			assertions: func(t *testing.T, _ arnetworkal.Connection, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "error establishing d2c connection")
				assert.Contains(t, err.Error(), "bat")
			},
		},

		{
			name: "establishing a connection succeeds",
			negotiateConnectionBehavior: func(net.IP, int, int) (int, error) {
				return 12345, nil
			},
			establishC2DConnectionBehavior: func(net.IP, int) (*net.UDPConn, error) {
				return nil, nil
			},
			establishD2CConnectionBehavior: func(int) (*net.UDPConn, error) {
				return nil, nil
			},
			assertions: func(t *testing.T, _ arnetworkal.Connection, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			negotiateConnection = defaultNegotiateConnection
			if testCase.negotiateConnectionBehavior != nil {
				negotiateConnection = testCase.negotiateConnectionBehavior
			}
			establishC2DConnection = defaultEstablishC2DConnection
			if testCase.establishC2DConnectionBehavior != nil {
				establishC2DConnection = testCase.establishC2DConnectionBehavior
			}
			establishD2CConnection = defaultEstablishD2CConnection
			if testCase.establishD2CConnectionBehavior != nil {
				establishD2CConnection = testCase.establishD2CConnectionBehavior
			}
			conn, err := NewConnection()
			testCase.assertions(t, conn, err)
			if conn != nil {
				conn.Close()
			}
		})
	}
}

func TestCannotConnectToNegotiate(t *testing.T) {
	// Pick a port where we're sure nothing else is listening
	discoveryPort, err := freeport.GetFreePort()
	require.NoError(t, err)
	_, err = defaultNegotiateConnection(
		net.ParseIP("127.0.0.1"),
		discoveryPort,
		12345, // Dummy port number-- we'll never connect to this, so it's ok
	)
	require.Error(t, err)
}

func TestRefusedConnectionNegotiation(t *testing.T) {
	discoveryPort, err := freeport.GetFreePort()
	require.NoError(t, err)
	// Dummy port number-- we'll never connect to this, so it's ok
	c2dPort := 12345
	// Mock out the device's connection negotiation. Run it in its on goroutine so
	// we can move on to trying to talk to it.
	listeningCh := make(chan struct{})
	go func() {
		// nolint: vetshadow
		listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: discoveryPort})
		require.NoError(t, err)
		defer listener.Close()
		close(listeningCh) // Signal the test to continue
		// Wait for a connection
		conn, err := listener.AcceptTCP()
		require.NoError(t, err)
		defer conn.Close()
		// Wait for the request
		data, err := bufio.NewReader(conn).ReadBytes(0x00)
		require.NoError(t, err)
		var negReq connectionNegotiationRequest
		err = json.Unmarshal(data[:len(data)-1], &negReq)
		require.NoError(t, err)
		// Send a response
		jsonBytes, err := json.Marshal(
			connectionNegotiationResponse{
				Status:  1, // Non-zero == connection refused
				C2DPort: c2dPort,
			},
		)
		require.NoError(t, err)
		jsonBytes = append(jsonBytes, 0x00)
		_, err = conn.Write(jsonBytes)
		require.NoError(t, err)
	}()

	// Block until mock device's connection negotiation server is listening, or
	// give up after 2 seconds. This prevents us from racing to connect to a
	// server that isn't listening yet. The timeout is to avoid the possibility of
	// blocking indefinitely if something goes wrong.
	select {
	case <-listeningCh:
	case <-time.After(2 * time.Second):
		require.Fail(
			t,
			"timed out waiting for mock device connection negotiation server to "+
				"start listening",
		)
	}

	_, err = defaultNegotiateConnection(
		net.ParseIP("127.0.0.1"),
		discoveryPort,
		// Dummy port number-- this is ok since the mock server won't do anything
		// with it
		54321,
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "refused")
}

func TestSuccessfulConnectionNegotiation(t *testing.T) {
	discoveryPort, err := freeport.GetFreePort()
	require.NoError(t, err)
	// Dummy port number-- we'll never connect to this, so it's ok
	c2dPort := 12345
	// Mock out the device's connection negotiation. Run it in its on goroutine so
	// we can move on to trying to talk to it.
	listeningCh := make(chan struct{})
	go func() {
		// nolint: vetshadow
		listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: discoveryPort})
		require.NoError(t, err)
		defer listener.Close()
		close(listeningCh) // Signal the test to continue
		// Wait for a connection
		conn, err := listener.AcceptTCP()
		require.NoError(t, err)
		defer conn.Close()
		// Wait for the request
		data, err := bufio.NewReader(conn).ReadBytes(0x00)
		require.NoError(t, err)
		var negReq connectionNegotiationRequest
		err = json.Unmarshal(data[:len(data)-1], &negReq)
		require.NoError(t, err)
		// Send a response
		jsonBytes, err := json.Marshal(
			connectionNegotiationResponse{
				Status:  0,
				C2DPort: c2dPort,
			},
		)
		require.NoError(t, err)
		jsonBytes = append(jsonBytes, 0x00)
		_, err = conn.Write(jsonBytes)
		require.NoError(t, err)
	}()

	// Block until mock device's connection negotiation server is listening, or
	// give up after 2 seconds. This prevents us from racing to connect to a
	// server that isn't listening yet. The timeout is to avoid the possibility of
	// blocking indefinitely if something goes wrong.
	select {
	case <-listeningCh:
	case <-time.After(2 * time.Second):
		require.Fail(
			t,
			"timed out waiting for mock device connection negotiation server to "+
				"start listening",
		)
	}

	negotiatedC2DPort, err := defaultNegotiateConnection(
		net.ParseIP("127.0.0.1"),
		discoveryPort,
		// Dummy port number-- this is ok since the mock server won't do anything
		// with it
		54321,
	)
	require.NoError(t, err)
	require.Equal(t, c2dPort, negotiatedC2DPort)
}

func TestC2DConnection(t *testing.T) {
	c2dPort, err := freeport.GetFreePort() // nolint: vetshadow
	require.NoError(t, err)
	conn := &connection{}
	conn.c2dConn, err = defaultEstablishC2DConnection(
		net.ParseIP("127.0.0.1"),
		c2dPort,
	)
	require.NoError(t, err)
	defer conn.closeC2DConnection()
	// Override frame encoding scheme to keep things simple-- we'll always
	// send "foo"
	conn.encodeFrame = func(arnetworkal.Frame) ([]byte, error) {
		return []byte("foo"), nil
	}

	// Set up a mock device c2d connection as a destination for c2d traffic
	mockDeviceC2DConn, err := net.ListenUDP(
		"udp",
		&net.UDPAddr{Port: conn.c2dConn.RemoteAddr().(*net.UDPAddr).Port},
	)
	require.NoError(t, err)
	defer mockDeviceC2DConn.Close()

	// Send some data to the mock device
	err = conn.Send(arnetworkal.Frame{})
	require.NoError(t, err)

	// Verify the mock device received the data
	packet := make([]byte, maxUDPDataBytes)
	bytesRead, _, err := mockDeviceC2DConn.ReadFromUDP(packet)
	require.NoError(t, err)
	require.Equal(t, "foo", string(packet[:bytesRead]))
}

func TestD2CConnection(t *testing.T) {
	d2cPort, err := freeport.GetFreePort() // nolint: vetshadows
	require.NoError(t, err)
	conn := &connection{}
	conn.d2cConn, err = defaultEstablishD2CConnection(d2cPort)
	require.NoError(t, err)
	defer conn.closeD2CConnection()
	// Override packet decoding to make some assertions
	conn.decodePacket = func(packet []byte) ([]arnetworkal.Frame, error) {
		// Expect to receive "bar"
		require.Equal(t, "bar", string(packet))
		return nil, nil
	}

	// Set up a mock device d2c connection as source of d2c traffic
	mockDeviceD2CConn, err := net.DialUDP(
		"udp",
		nil,
		&net.UDPAddr{Port: conn.d2cConn.LocalAddr().(*net.UDPAddr).Port},
	)
	require.NoError(t, err)
	defer mockDeviceD2CConn.Close()

	// Make the mock device send some data
	_, err = mockDeviceD2CConn.Write([]byte("bar"))
	require.NoError(t, err)

	// Receive some data from the mock device
	_, err = conn.Receive()
	require.NoError(t, err)
}
