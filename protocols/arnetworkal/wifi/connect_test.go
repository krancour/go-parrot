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

func TestConnect(t *testing.T) {
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
		assertions                     func(
			*testing.T,
			arnetworkal.FrameSender,
			arnetworkal.FrameReceiver,
			error,
		)
	}{

		{
			name: "connection negotiation fails",
			negotiateConnectionBehavior: func(net.IP, int, int) (int, error) {
				return 0, errors.New("foo")
			},
			assertions: func(
				t *testing.T,
				_ arnetworkal.FrameSender,
				_ arnetworkal.FrameReceiver,
				err error,
			) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "connection negotiation failed")
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
			assertions: func(
				t *testing.T,
				_ arnetworkal.FrameSender,
				_ arnetworkal.FrameReceiver,
				err error,
			) {
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
			assertions: func(
				t *testing.T,
				_ arnetworkal.FrameSender,
				_ arnetworkal.FrameReceiver,
				err error,
			) {
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
			assertions: func(
				t *testing.T,
				_ arnetworkal.FrameSender,
				_ arnetworkal.FrameReceiver,
				err error,
			) {
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
			frameSender, frameReceiver, err := Connect()
			testCase.assertions(t, frameSender, frameReceiver, err)
			if frameSender != nil {
				frameSender.Close()
			}
			if frameReceiver != nil {
				frameReceiver.Close()
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
