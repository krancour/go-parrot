package ardrone3

import (
	"bytes"
	"encoding/binary"
)

type pcmd struct {
	Flag  int
	Roll  int
	Pitch int
	Yaw   int
	Gaz   int
	Psi   float32
}

func validatePitch(val int) int {
	if val > 100 {
		return 100
	} else if val < 0 {
		return 0
	}
	return val
}

func (b *BebopClient) generatePcmd() *bytes.Buffer {
	//
	// ARCOMMANDS_Generator_GenerateARDrone3PilotingPCMD
	//
	// uint8 - flag Boolean flag to activate roll/pitch movement
	// int8  - roll Roll consign for the drone [-100;100]
	// int8  - pitch Pitch consign for the drone [-100;100]
	// int8  - yaw Yaw consign for the drone [-100;100]
	// int8  - gaz Gaz consign for the drone [-100;100]
	// float - psi [NOT USED] - Magnetic north heading of the
	//         controlling device (deg) [-180;180]
	//

	cmd := &bytes.Buffer{}
	tmp := &bytes.Buffer{}

	cmd.WriteByte(projectARDrone3)
	cmd.WriteByte(classPiloting)

	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, uint16(cmdPCMD))
	cmd.Write(tmp.Bytes())

	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, uint8(b.Pcmd.Flag))
	cmd.Write(tmp.Bytes())

	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, int8(b.Pcmd.Roll))
	cmd.Write(tmp.Bytes())

	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, int8(b.Pcmd.Pitch))
	cmd.Write(tmp.Bytes())

	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, int8(b.Pcmd.Yaw))
	cmd.Write(tmp.Bytes())

	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, int8(b.Pcmd.Gaz))
	cmd.Write(tmp.Bytes())

	tmp = &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, uint32(b.Pcmd.Psi))
	cmd.Write(tmp.Bytes())

	return b.networkFrameGenerator(cmd, frameTypeData, bufferNonack)
}

func (b *BebopClient) TakeOff() error {
	cmd := &bytes.Buffer{}
	cmd.WriteByte(projectARDrone3)
	cmd.WriteByte(classPiloting)
	tmp := &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, uint16(cmdTakeOff))
	cmd.Write(tmp.Bytes())
	_, err := b.write(b.networkFrameGenerator(cmd, frameTypeData, bufferNonack).Bytes())
	return err
}

func (b *BebopClient) Land() error {
	cmd := &bytes.Buffer{}
	cmd.WriteByte(projectARDrone3)
	cmd.WriteByte(classPiloting)
	tmp := &bytes.Buffer{}
	binary.Write(tmp, binary.LittleEndian, uint16(cmdLanding))
	cmd.Write(tmp.Bytes())
	_, err := b.write(b.networkFrameGenerator(cmd, frameTypeData, bufferNonack).Bytes())
	return err
}

func (b *BebopClient) Up(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Gaz = validatePitch(val)
	return nil
}

func (b *BebopClient) Down(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Gaz = validatePitch(val) * -1
	return nil
}

func (b *BebopClient) Forward(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Pitch = validatePitch(val)
	return nil
}

func (b *BebopClient) Backward(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Pitch = validatePitch(val) * -1
	return nil
}

func (b *BebopClient) Right(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Roll = validatePitch(val)
	return nil
}

func (b *BebopClient) Left(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Roll = validatePitch(val) * -1
	return nil
}

func (b *BebopClient) Clockwise(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Yaw = validatePitch(val)
	return nil
}

func (b *BebopClient) CounterClockwise(val int) error {
	b.Pcmd.Flag = 1
	b.Pcmd.Yaw = validatePitch(val) * -1
	return nil
}

func (b *BebopClient) Stop() error {
	b.Pcmd = pcmd{
		Flag:  0,
		Roll:  0,
		Pitch: 0,
		Yaw:   0,
		Gaz:   0,
		Psi:   0,
	}
	return nil
}
