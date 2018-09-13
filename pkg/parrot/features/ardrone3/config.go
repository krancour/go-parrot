package ardrone3

import (
	"bytes"
	"encoding/binary"
)

func (b *BebopClient) FlatTrim() error {
	cmd := &bytes.Buffer{}
	cmd.WriteByte(projectARDrone3)
	cmd.WriteByte(classPiloting)
	tmp := &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, uint16(cmdFlatTrim))
	cmd.Write(tmp.Bytes())
	_, err := b.write(b.networkFrameGenerator(cmd, frameTypeData, bufferNonack).Bytes())
	return err
}

func (b *BebopClient) HullProtection(protect bool) error {
	cmd := &bytes.Buffer{}
	cmd.WriteByte(projectARDrone3)
	cmd.WriteByte(classSpeedSettings)
	tmp := &bytes.Buffer{}
	binary.Write(tmp,
		binary.LittleEndian,
		uint16(cmdHullProtection),
	)
	cmd.Write(tmp.Bytes())
	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, bool2int8(protect))
	cmd.Write(tmp.Bytes())
	_, err := b.write(b.networkFrameGenerator(cmd, frameTypeData, bufferNonack).Bytes())
	return err
}

func (b *BebopClient) Outdoor(outdoor bool) error {
	cmd := &bytes.Buffer{}
	cmd.WriteByte(projectARDrone3)
	cmd.WriteByte(classSpeedSettings)
	tmp := &bytes.Buffer{}
	binary.Write(tmp,
		binary.LittleEndian,
		uint16(cmdOutdoor),
	)
	cmd.Write(tmp.Bytes())
	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, bool2int8(outdoor))
	cmd.Write(tmp.Bytes())
	_, err := b.write(b.networkFrameGenerator(cmd, frameTypeData, bufferNonack).Bytes())
	return err
}
