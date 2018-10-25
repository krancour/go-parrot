package wifi

import (
	"bufio"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

// TODO: Break this into smaller, more focused test cases
func TestConnection(t *testing.T) {
	// Override the IP address and ports
	deviceIP = net.ParseIP("127.0.0.1")

	// Pick an available port and override the default discovery port.
	// (This is also the port used for connection negotiation.)
	var err error
	discoveryPort, err = freeport.GetFreePort()
	require.NoError(t, err)

	// Mock out the device's connection negotiation. Run it in its on goroutine so
	// we can move on to trying to talk to it.
	listeningCh := make(chan struct{})
	go func() {
		c2dPort, err := freeport.GetFreePort() // nolint: vetshadow
		require.NoError(t, err)
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
	// give up after 5 seconds. This prevents us from racing to connect to a
	// server that isn't listening yet. The timeout is to avoid the possibility of
	// blocking indefinitely if something goes wrong.
	select {
	case <-listeningCh:
	case <-time.After(5 * time.Second):
		require.Fail(
			t,
			"timed out waiting for mock device connection negotiation server to "+
				"start listening",
		)
	}

	// Create a UDP/IP based ARNetworkAL connection
	iConn, err := NewConnection()
	require.NoError(t, err)
	defer iConn.Close()
	conn, ok := iConn.(*connection)
	require.True(t, ok)

	// Override frame encoding scheme to keep things simple-- we'll always
	// send "foo"
	conn.encodeFrame = func(arnetworkal.Frame) ([]byte, error) {
		return []byte("foo"), nil
	}

	// Override packet decoding to make some assertions
	conn.decodePacket = func(packet []byte) ([]arnetworkal.Frame, error) {
		// Expect to receive "bar"
		require.Equal(t, "bar", string(packet))
		return nil, nil
	}

	// Set up a mock device c2d connection as a destination for c2d traffic
	mockDeviceC2DConn, err := net.ListenUDP(
		"udp",
		&net.UDPAddr{Port: conn.c2dConn.RemoteAddr().(*net.UDPAddr).Port},
	)
	require.NoError(t, err)
	defer mockDeviceC2DConn.Close()

	// Set up a mock device d2c connection as source of d2c traffic
	mockDeviceD2CConn, err := net.DialUDP(
		"udp",
		nil,
		&net.UDPAddr{Port: conn.d2cConn.LocalAddr().(*net.UDPAddr).Port},
	)
	require.NoError(t, err)
	defer mockDeviceD2CConn.Close()

	// Send some data to the mock device
	err = conn.Send(arnetworkal.Frame{})
	require.NoError(t, err)

	// Verify the mock device received the data
	packet := make([]byte, maxUDPDataBytes)
	bytesRead, _, err := mockDeviceC2DConn.ReadFromUDP(packet)
	require.NoError(t, err)
	require.Equal(t, "foo", string(packet[:bytesRead]))

	// Make the mock device send some data
	_, err = mockDeviceD2CConn.Write([]byte("bar"))
	require.NoError(t, err)

	// Receive some data from the mock device
	_, err = conn.Receive()
	require.NoError(t, err)
}
