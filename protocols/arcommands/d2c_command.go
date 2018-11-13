package arcommands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// D2CCommand ...
// TODO: Document this
type D2CCommand interface {
	execute(data []byte) error
}

type d2cCommand struct {
	name        string
	argTemplate []interface{}
	callback    func(args []interface{}) error
}

// NewD2CCommand ...
// TODO: Document this
func NewD2CCommand(
	name string,
	argTemplate []interface{},
	callback func(args []interface{}) error,
) D2CCommand {
	return &d2cCommand{
		name:        name,
		argTemplate: argTemplate,
		callback:    callback,
	}
}

func (d *d2cCommand) execute(data []byte) error {
	// TODO: Make sure this is full of zero values!
	args := d.argTemplate
	if err := decodeArgs(data, args); err != nil {
		return err
	}
	err := d.callback(args)
	return err
}

func decodeArgs(data []byte, args []interface{}) error {
	buf := bytes.NewBuffer(data[4:len(data)])
	var err error
	for i, argIface := range args {
		switch arg := argIface.(type) {
		case uint8:
			err = binary.Read(buf, binary.LittleEndian, &arg)
			args[i] = arg
		case int8:
			err = binary.Read(buf, binary.LittleEndian, &arg)
			args[i] = arg
		case uint16:
			err = binary.Read(buf, binary.LittleEndian, &arg)
			args[i] = arg
		case int16:
			err = binary.Read(buf, binary.LittleEndian, &arg)
			args[i] = arg
		case uint32:
			err = binary.Read(buf, binary.LittleEndian, &arg)
			args[i] = arg
		case int32:
			err = binary.Read(buf, binary.LittleEndian, &arg)
			args[i] = arg
		case uint64:
			err = binary.Read(buf, binary.LittleEndian, &arg)
			args[i] = arg
		case int64:
			err = binary.Read(buf, binary.LittleEndian, &arg)
			args[i] = arg
		case float32:
			err = binary.Read(buf, binary.LittleEndian, &arg)
			args[i] = arg
		case float64:
			err = binary.Read(buf, binary.LittleEndian, &arg)
			args[i] = arg
		case string:
			var bytes []byte
			bytes, err = buf.ReadBytes(0x00)
			bytes = bytes[0 : len(bytes)-1]
			args[i] = string(bytes)
		default:
			err = fmt.Errorf("unknown type: %s", reflect.TypeOf(argIface))
		}
		if err != nil {
			return err
		}
	}
	return nil
}
