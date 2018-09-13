package wifi

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/krancour/drone-examples/pkg/parrot/protocols/arnetworkal"
)

// headerBytesLength is the combined length of all ARNetworkAL frame headers
// in bytes.
const headerBytesLength = 7

func defaultEncodeFrame(frame arnetworkal.Frame) []byte {
	var packetBuf bytes.Buffer
	packetBuf.WriteByte(byte(frame.Type)) // 1 byte
	packetBuf.WriteByte(frame.ID)         // 1 byte
	packetBuf.WriteByte(frame.Seq)        // 1 byte
	var sizeBuf bytes.Buffer
	binary.Write(
		&sizeBuf,
		binary.LittleEndian,
		uint32(headerBytesLength+len(frame.Data)),
	)
	packetBuf.Write(sizeBuf.Bytes()) // 4 bytes
	packetBuf.Write(frame.Data)
	return packetBuf.Bytes()
}

func defaultDecodePacket(packet []byte) ([]arnetworkal.Frame, error) {
	data := packet
	frames := []arnetworkal.Frame{}
	for {
		if len(data) == 0 {
			return frames, nil
		}
		if len(data) < headerBytesLength {
			// We are clearly dealing with a malformed packet. We can't trust
			// ANY of these frames. Discard them all and return an error.
			return nil, fmt.Errorf("error decoding malformed packet")
		}
		frame := arnetworkal.Frame{
			Type: arnetworkal.FrameType(data[0]), // 1 byte
			ID:   data[1],                        // 1 byte
			Seq:  data[2],                        // 1 byte
			Data: []byte{},
		}
		var frameSize uint32 // 4 bytes
		binary.Read(
			bytes.NewReader(
				data[3:7],
			),
			binary.LittleEndian,
			&frameSize,
		)
		if uint32(len(data)) < frameSize {
			// We are clearly dealing with a malformed packet. We can't trust
			// ANY of these frames. Discard them all and return an error.
			return nil, fmt.Errorf("error decoding malformed packet")
		}
		frame.Data = data[7:frameSize]
		frames = append(frames, frame)
		data = data[frameSize:]
	}
}
