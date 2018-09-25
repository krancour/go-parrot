package wifi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	// Override the IP address and ports
	deviceIP = "127.0.0.1"

	// Pick an available port and override the default discovery port.
	// (This is also the port used for connection negotiation.)
	var err error
	discoveryPort, err = freeport.GetFreePort()
	assert.Nil(t, err)

	// Create a dummy server that implements connection negotiation. Run it in
	// its on goroutine so we can move on to trying to talk to it.
	listeningCh := make(chan struct{})
	go func() {
		c2dPort, err := freeport.GetFreePort()
		assert.Nil(t, err)
		negConnAddr, err := net.ResolveTCPAddr(
			"tcp",
			fmt.Sprintf(":%d", discoveryPort),
		)
		assert.Nil(t, err)
		negListener, err := net.ListenTCP("tcp", negConnAddr)
		defer negListener.Close()
		assert.Nil(t, err)
		close(listeningCh) // Signal the test to continue
		// Wait for a connection
		negConn, err := negListener.AcceptTCP()
		assert.Nil(t, err)
		defer negConn.Close()
		// Wait for the request
		data, err := bufio.NewReader(negConn).ReadBytes(0x00)
		assert.Nil(t, err)
		var negReq connectionNegotiationRequest
		err = json.Unmarshal(data[:len(data)-1], &negReq)
		assert.Nil(t, err)
		// Send a response
		jsonBytes, err := json.Marshal(
			connectionNegotiationResponse{
				Status:  0,
				C2DPort: c2dPort,
			},
		)
		jsonBytes = append(jsonBytes, 0x00)
		_, err = negConn.Write(jsonBytes)
		assert.Nil(t, err)
	}()

	// Block until the test server is listening, or give up after 5 seconds.
	// This prevents us from racing to connect to a server that isn't listening
	// yet. The timeout is to avoid the possibility of blocking indefinitely if
	// something goes wrong.
	select {
	case <-listeningCh:
	case <-time.After(5 * time.Second):
		require.Fail(t, "timed out waiting for test server to start listening")
	}

	// Create a UDP/IP based ARNetworkAL connection
	iConn, err := NewConnection()
	assert.Nil(t, err)
	defer iConn.Close()
	conn, ok := iConn.(*connection)
	assert.True(t, ok)

	// Override frame encoding scheme to keep things simple-- we'll always
	// send "foo"
	conn.encodeFrame = func(arnetworkal.Frame) []byte {
		return []byte("foo")
	}

	// Override packet decoding to make some assertions
	conn.decodePacket = func(packet []byte) ([]arnetworkal.Frame, error) {
		// Expect to receive "bar"
		assert.Equal(t, "bar", string(packet))
		return nil, nil
	}

	// Set up a test server we can send to
	testServerListenAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf(":%d", conn.c2dPort),
	)
	assert.Nil(t, err)
	testServerListenConn, err := net.ListenUDP("udp", testServerListenAddr)
	assert.Nil(t, err)
	defer testServerListenConn.Close()

	// Set up a test client we can receive from
	testClientWriteAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf("%s:%d", deviceIP, conn.d2cPort),
	)
	assert.Nil(t, err)
	testClientWriteConn, err := net.DialUDP("udp", nil, testClientWriteAddr)
	assert.Nil(t, err)
	defer testClientWriteConn.Close()

	// Use the connection to send some data to the test server
	err = conn.Send(arnetworkal.Frame{})
	assert.Nil(t, err)

	// Verify the test server received the data
	packet := make([]byte, maxUDPDataBytes)
	bytesRead, _, err := testServerListenConn.ReadFromUDP(packet)
	assert.Nil(t, err)
	assert.Equal(t, "foo", string(packet[:bytesRead]))

	// Use the test client to send some data for the connection to receive
	_, err = testClientWriteConn.Write([]byte("bar"))
	assert.Nil(t, err)

	// Use the connection to receive some data from the test client
	_, err = conn.Receive()
	assert.Nil(t, err)
}
