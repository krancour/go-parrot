package arcommands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// D2CCommand ...
// TODO: Document this
type D2CCommand interface {
	ID() uint16
	Name() string
	execute(data []byte) error
}

type d2cCommand struct {
	id          uint16
	name        string
	argTemplate []interface{}
	callback    func(args []interface{}) error
}

// NewD2CCommand ...
// TODO: Document this
func NewD2CCommand(
	id uint16,
	name string,
	argTemplate []interface{},
	callback func(args []interface{}) error,
) D2CCommand {
	return &d2cCommand{
		id:          id,
		name:        name,
		argTemplate: argTemplate,
		callback:    callback,
	}
}

func (d *d2cCommand) ID() uint16 {
	return d.id
}

func (d *d2cCommand) Name() string {
	return d.name
}

func (d *d2cCommand) execute(data []byte) error {
	// Super important-- make a COPY of the argument template!
	args := make([]interface{}, len(d.argTemplate))
	copy(args, d.argTemplate)
	if err := decodeArgs(data, args); err != nil {
		return errors.Wrap(err, "error decoding command arguments")
	}
	if err := d.callback(args); err != nil {
		return errors.Wrap(err, "error executing command")
	}
	return nil
}

func decodeArgs(data []byte, args []interface{}) error {
	buf := bytes.NewBuffer(data[4:])
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
			return errors.Wrapf(
				err,
				"error decoding command arguments; data: %v: args: %v",
				data,
				args,
			)
		}
	}
	return nil
}
