package wifi

import (
	"bytes"
	"encoding/binary"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// headerBytesLength is the combined length of all ARNetworkAL frame headers
// in bytes.
const headerBytesLength = 7

func defaultEncodeFrame(frame arnetworkal.Frame) ([]byte, error) {
	log := log.WithField(
		"uuid", frame.UUID,
	).WithField(
		"buffer", frame.ID,
	).WithField(
		"type", frame.Type,
	).WithField(
		"seq", frame.Seq,
	)
	log.Debug("encoding arnetworkal frame as datagram")
	var datagramBuf bytes.Buffer
	datagramBuf.WriteByte(byte(frame.Type)) // 1 byte
	datagramBuf.WriteByte(frame.ID)         // 1 byte
	datagramBuf.WriteByte(frame.Seq)        // 1 byte
	var sizeBuf bytes.Buffer
	if err := binary.Write(
		&sizeBuf,
		binary.LittleEndian,
		uint32(headerBytesLength+len(frame.Data)),
	); err != nil {
		return nil,
			errors.Wrap(err, "error encoding arnetworkal frame as datagram")
	}
	datagramBuf.Write(sizeBuf.Bytes()) // 4 bytes
	datagramBuf.Write(frame.Data)
	log.Debug("encoded arnetworkal frame as datagram")
	return datagramBuf.Bytes(), nil
}

func defaultDecodeDatagram(datagram []byte) ([]arnetworkal.Frame, error) {
	log.Debug("decoding datagram")
	data := datagram
	frames := []arnetworkal.Frame{}
	for {
		if len(data) == 0 {
			log.WithField(
				"frameCount", len(frames),
			).Debug("extracted arnetworkal frames from datagram")
			return frames, nil
		}
		if len(data) < headerBytesLength {
			// We are clearly dealing with a malformed datagram. We can't trust
			// ANY of these frames. Discard them all and return an error.
			return nil, errors.New("error decoding malformed datagram")
		}
		var frameUUID string
		if log.GetLevel() == log.DebugLevel {
			frameUUID = uuid.NewV4().String()
		}
		frame := arnetworkal.Frame{
			UUID: frameUUID,
			Type: arnetworkal.FrameType(data[0]), // 1 byte
			ID:   data[1],                        // 1 byte
			Seq:  data[2],                        // 1 byte
			Data: []byte{},
		}
		var frameSize uint32 // 4 bytes
		if err := binary.Read(
			bytes.NewReader(
				data[3:7],
			),
			binary.LittleEndian,
			&frameSize,
		); err != nil {
			return nil, errors.Wrap(
				err,
				"error determining arnetworkal frame data length",
			)
		}
		if uint32(len(data)) < frameSize {
			// We are clearly dealing with a malformed datagram. We can't trust
			// ANY of these frames. Discard them all and return an error.
			return nil, errors.New("error decoding malformed datagram")
		}
		frame.Data = data[7:frameSize]
		log.WithField(
			"uuid", frame.UUID,
		).WithField(
			"buffer", frame.ID,
		).WithField(
			"type", frame.Type,
		).WithField(
			"seq", frame.Seq,
		).Debug("extracted arnetworkal frame from datagram")
		frames = append(frames, frame)
		data = data[frameSize:]
	}
}
