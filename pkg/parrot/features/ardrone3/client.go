package ardrone3

import (
	"bytes"
	"encoding/binary"
	"net"
)

type BebopClient struct {
	NavData               map[string]string
	Pcmd                  pcmd
	c2dClient             *net.UDPConn
	d2cClient             *net.UDPConn
	discoveryClient       *net.TCPConn
	networkFrameGenerator func(*bytes.Buffer, byte, byte) *bytes.Buffer
	writeChan             chan []byte
}

func New() *BebopClient {
	return &BebopClient{
		NavData:               make(map[string]string),
		networkFrameGenerator: networkFrameGenerator(),
		Pcmd: pcmd{
			Flag:  0,
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
			Gaz:   0,
			Psi:   0,
		},
		writeChan: make(chan []byte),
	}
}

func (b *BebopClient) GenerateAllStates() error {
	cmd := &bytes.Buffer{}
	cmd.WriteByte(projectCommon)
	cmd.WriteByte(classCommon)
	tmp := &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, uint16(cmdAllStates))
	cmd.Write(tmp.Bytes())
	_, err := b.write(b.networkFrameGenerator(cmd, frameTypeData, bufferNonack).Bytes())
	return err
}
