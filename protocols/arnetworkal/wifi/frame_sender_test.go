package wifi

import (
	"net"
	"testing"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

func TestSendFrame(t *testing.T) {
	c2dPort, err := freeport.GetFreePort() // nolint: vetshadow
	require.NoError(t, err)
	frameSender := &frameSender{}
	frameSender.conn, err = defaultEstablishC2DConnection(
		net.ParseIP("127.0.0.1"),
		c2dPort,
	)
	require.NoError(t, err)
	defer frameSender.Close()
	// Override frame encoding scheme to keep things simple-- we'll always
	// send "foo"
	frameSender.encodeFrame = func(arnetworkal.Frame) ([]byte, error) {
		return []byte("foo"), nil
	}

	// Set up a mock device c2d connection as a destination for c2d traffic
	mockDeviceC2DConn, err := net.ListenUDP(
		"udp",
		&net.UDPAddr{Port: frameSender.conn.RemoteAddr().(*net.UDPAddr).Port},
	)
	require.NoError(t, err)
	defer mockDeviceC2DConn.Close()

	// Send some data to the mock device
	err = frameSender.Send(arnetworkal.Frame{})
	require.NoError(t, err)

	// Verify the mock device received the data
	datagram := make([]byte, maxUDPDataBytes)
	bytesRead, _, err := mockDeviceC2DConn.ReadFromUDP(datagram)
	require.NoError(t, err)
	require.Equal(t, "foo", string(datagram[:bytesRead]))
}
