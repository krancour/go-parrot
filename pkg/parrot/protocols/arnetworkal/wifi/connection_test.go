package wifi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"github.com/krancour/drone-examples/pkg/parrot/protocols/arnetworkal"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
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

	// Override data decoding to make some assertions
	conn.decodeData = func(data []byte) ([]arnetworkal.Frame, error) {
		// Expect to receive "bar"
		assert.Equal(t, "bar", string(data))
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
	data := make([]byte, maxUDPDataBytes)
	bytesRead, _, err := testServerListenConn.ReadFromUDP(data)
	assert.Nil(t, err)
	assert.Equal(t, "foo", string(data[:bytesRead]))

	// Use the test client to send some data for the connection to receive
	_, err = testClientWriteConn.Write([]byte("bar"))
	assert.Nil(t, err)

	// Use the connection to receive some data from the test client
	_, err = conn.Receive()
	assert.Nil(t, err)
}
