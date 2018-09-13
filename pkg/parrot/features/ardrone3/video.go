package ardrone3

import (
	"bytes"
	"encoding/binary"
)

func (b *BebopClient) StartRecording() error {
	buf := b.videoRecord(argVideoRecordStart)
	b.write(b.networkFrameGenerator(buf, frameTypeData, bufferNonack).Bytes())
	return nil
}

func (b *BebopClient) StopRecording() error {
	buf := b.videoRecord(argVideoRecordStop)
	b.write(b.networkFrameGenerator(buf, frameTypeData, bufferNonack).Bytes())
	return nil
}

func (b *BebopClient) videoRecord(state byte) *bytes.Buffer {
	cmd := &bytes.Buffer{}
	cmd.WriteByte(projectARDrone3)
	cmd.WriteByte(classMediaRecord)
	tmp := &bytes.Buffer{}
	binary.Write(tmp,
		binary.LittleEndian,
		uint16(cmdVideo),
	)
	cmd.Write(tmp.Bytes())
	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, uint32(state))
	cmd.Write(tmp.Bytes())
	cmd.WriteByte(0)
	return cmd
}

func (b *BebopClient) VideoEnable(enable bool) error {
	cmd := &bytes.Buffer{}
	cmd.WriteByte(projectARDrone3)
	cmd.WriteByte(classMediaStreaming)
	tmp := &bytes.Buffer{}
	binary.Write(tmp,
		binary.LittleEndian,
		uint16(cmdVideoEnable),
	)
	cmd.Write(tmp.Bytes())
	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, bool2int8(enable))
	cmd.Write(tmp.Bytes())
	_, err := b.write(b.networkFrameGenerator(cmd, frameTypeData, bufferNonack).Bytes())
	return err
}

func (b *BebopClient) VideoStreamMode(mode int8) error {
	cmd := &bytes.Buffer{}
	cmd.WriteByte(projectARDrone3)
	cmd.WriteByte(classMediaStreaming)
	tmp := &bytes.Buffer{}
	binary.Write(tmp,
		binary.LittleEndian,
		uint16(cmdVideoStreamMode),
	)
	cmd.Write(tmp.Bytes())
	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, mode)
	cmd.Write(tmp.Bytes())
	_, err := b.write(b.networkFrameGenerator(cmd, frameTypeData, bufferNonack).Bytes())
	return err
}
