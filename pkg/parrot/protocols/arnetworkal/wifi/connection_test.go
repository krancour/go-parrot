package wifi

import (
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
	var err error
	c2dPort, err = freeport.GetFreePort()
	assert.Nil(t, err)
	d2cPort, err = freeport.GetFreePort()
	assert.Nil(t, err)

	// Set up a test server we can send to
	testServerListenAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf(":%d", c2dPort),
	)
	assert.Nil(t, err)
	testServerListenConn, err := net.ListenUDP("udp", testServerListenAddr)
	assert.Nil(t, err)
	defer testServerListenConn.Close()

	// Set up a test client we can receive from
	testClientWriteAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf("%s:%d", deviceIP, d2cPort),
	)
	assert.Nil(t, err)
	testClientWriteConn, err := net.DialUDP("udp", nil, testClientWriteAddr)
	assert.Nil(t, err)
	defer testClientWriteConn.Close()

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
