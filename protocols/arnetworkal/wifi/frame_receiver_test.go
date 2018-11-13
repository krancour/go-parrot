package wifi

import (
	"net"
	"testing"

	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

func TestReceiveFrame(t *testing.T) {
	d2cPort, err := freeport.GetFreePort() // nolint: vetshadows
	require.NoError(t, err)
	frameReceiver := &frameReceiver{
		datagramBuffer: make([]byte, maxUDPDataBytes),
	}
	frameReceiver.conn, err = defaultEstablishD2CConnection(d2cPort)
	require.NoError(t, err)
	defer frameReceiver.Close()
	// Override datagram decoding to make some assertions
	frameReceiver.decodeDatagram =
		func(datagram []byte) ([]arnetworkal.Frame, error) {
			// Expect to receive "bar"
			require.Equal(t, "bar", string(datagram))
			return nil, nil
		}

	// Set up a mock device d2c connection as source of d2c traffic
	mockDeviceD2CConn, err := net.DialUDP(
		"udp",
		nil,
		&net.UDPAddr{Port: frameReceiver.conn.LocalAddr().(*net.UDPAddr).Port},
	)
	require.NoError(t, err)
	defer mockDeviceD2CConn.Close()

	// Make the mock device send some data
	_, err = mockDeviceD2CConn.Write([]byte("bar"))
	require.NoError(t, err)

	// Receive some data from the mock device
	_, err = frameReceiver.Receive()
	require.NoError(t, err)
}
