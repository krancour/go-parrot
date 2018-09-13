package ardrone3

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type networkFrame struct {
	Type int
	ID   int
	Seq  int
	Size int
	Data []byte
}

func newNetworkFrame(buf []byte) networkFrame {
	var size uint32
	binary.Read(bytes.NewReader(buf[3:7]), binary.LittleEndian, &size)
	return networkFrame{
		Type: int(buf[0]),
		ID:   int(buf[1]),
		Seq:  int(buf[2]),
		Size: int(size),
		Data: buf[7:size],
	}
}

func networkFrameGenerator() func(*bytes.Buffer, byte, byte) *bytes.Buffer {
	//
	// ARNETWORKAL_Frame_t
	//
	// uint8  type  - frame type ARNETWORK_FRAME_TYPE
	// uint8  id    - identifier of the buffer sending the frame
	// uint8  seq   - sequence number of the frame
	// uint32 size  - size of the frame
	//

	// each frame id has it's own sequence number
	seq := make(map[byte]byte)

	hlen := 7 // size of ARNETWORKAL_Frame_t header

	return func(cmd *bytes.Buffer, frameType byte, id byte) *bytes.Buffer {
		if _, ok := seq[id]; !ok {
			seq[id] = 0
		}

		seq[id]++

		if seq[id] > 255 {
			seq[id] = 0
		}

		ret := &bytes.Buffer{}
		ret.WriteByte(frameType)
		ret.WriteByte(id)
		ret.WriteByte(seq[id])

		size := &bytes.Buffer{}
		binary.Write(size, binary.LittleEndian, uint32(cmd.Len()+hlen))

		ret.Write(size.Bytes())
		ret.Write(cmd.Bytes())

		return ret
	}
}

func (b *BebopClient) write(buf []byte) (int, error) {
	b.writeChan <- buf
	return 0, nil
}

func (b *BebopClient) createAck(frame networkFrame) *bytes.Buffer {
	return b.networkFrameGenerator(bytes.NewBuffer([]byte{uint8(frame.Seq)}),
		frameTypeAck,
		byte(uint16(frame.ID)+(ARNETWORKAL_MANAGER_DEFAULT_ID_MAX/2)),
	)
}

func (b *BebopClient) createPong(frame networkFrame) *bytes.Buffer {
	return b.networkFrameGenerator(bytes.NewBuffer(frame.Data),
		frameTypeData,
		ARNETWORK_MANAGER_INTERNAL_BUFFER_ID_PONG,
	)
}

func (b *BebopClient) packetReceiver(buf []byte) {
	frame := newNetworkFrame(buf)
	if frame.Type == int(frameTypeDataWithAck) {
		ack := b.createAck(frame).Bytes()
		_, err := b.write(ack)

		if err != nil {
			fmt.Println("ARNETWORKAL_FRAME_TYPE_DATA_WITH_ACK", err)
		}
	}
	if frame.ID == int(ARNETWORK_MANAGER_INTERNAL_BUFFER_ID_PING) {
		pong := b.createPong(frame).Bytes()
		_, err := b.write(pong)
		if err != nil {
			fmt.Println("ARNETWORK_MANAGER_INTERNAL_BUFFER_ID_PING", err)
		}
	}
}
