package wifi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"github.com/krancour/drone-examples/protocols/arnetworkal"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	// Override the IP address and ports
	deviceIP = "127.0.0.1"

	// Pick an available port and override the default discovery port.
	// (This is also the port used for connection negotiation.)
	var err error
	discoveryPort, err = freeport.GetFreePort()
	require.Nil(t, err)
	// Create a dummy server that implements connection negotiation. Run it in
	// its on goroutine so we can move on to trying to talk to it.
	go func() {
		c2dPort, err := freeport.GetFreePort()
		require.Nil(t, err)
		negConnAddr, err := net.ResolveTCPAddr(
			"tcp",
			fmt.Sprintf(":%d", discoveryPort),
		)
		require.Nil(t, err)
		negListener, err := net.ListenTCP("tcp", negConnAddr)
		defer negListener.Close()
		require.Nil(t, err)
		// Wait for a connection
		negConn, err := negListener.Accept()
		require.Nil(t, err)
		defer negConn.Close()
		// Wait for the request
		data, err := bufio.NewReader(negConn).ReadBytes(0x00)
		require.Nil(t, err)
		var negReq connectionNegotiationRequest
		err = json.Unmarshal(data[:len(data)-1], &negReq)
		require.Nil(t, err)
		// Send a response
		jsonBytes, err := json.Marshal(
			connectionNegotiationResponse{
				Status:  0,
				C2DPort: c2dPort,
			},
		)
		jsonBytes = append(jsonBytes, 0x00)
		_, err = negConn.Write(jsonBytes)
		require.Nil(t, err)
	}()

	// Create a UDP/IP based ARNetworkAL connection
	iConn, err := NewConnection()
	require.Nil(t, err)
	defer iConn.Close()
	conn, ok := iConn.(*connection)
	require.True(t, ok)

	// Override frame encoding scheme to keep things simple-- we'll always
	// send "foo"
	conn.encodeFrame = func(arnetworkal.Frame) []byte {
		return []byte("foo")
	}

	// Override data decoding to make some assertions
	conn.decodeData = func(data []byte) ([]arnetworkal.Frame, error) {
		// Expect to receive "bar"
		require.Equal(t, "bar", string(data))
		return nil, nil
	}

	// Set up a test server we can send to
	testServerListenAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf(":%d", conn.c2dPort),
	)
	require.Nil(t, err)
	testServerListenConn, err := net.ListenUDP("udp", testServerListenAddr)
	require.Nil(t, err)
	defer testServerListenConn.Close()

	// Set up a test client we can receive from
	testClientWriteAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf("%s:%d", deviceIP, conn.d2cPort),
	)
	require.Nil(t, err)
	testClientWriteConn, err := net.DialUDP("udp", nil, testClientWriteAddr)
	require.Nil(t, err)
	defer testClientWriteConn.Close()

	// Use the connection to send some data to the test server
	err = conn.Send(arnetworkal.Frame{})
	require.Nil(t, err)

	// Verify the test server received the data
	data := make([]byte, maxUDPDataBytes)
	bytesRead, _, err := testServerListenConn.ReadFromUDP(data)
	require.Nil(t, err)
	require.Equal(t, "foo", string(data[:bytesRead]))

	// Use the test client to send some data for the connection to receive
	_, err = testClientWriteConn.Write([]byte("bar"))
	require.Nil(t, err)

	// Use the connection to receive some data from the test client
	_, err = conn.Receive()
	require.Nil(t, err)
}
