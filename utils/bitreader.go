package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
)

func GetData(b *bytes.Buffer, ptrs ...interface{}) (err error) {

	for i := range ptrs {
		switch data := ptrs[i].(type) {
		case *uint8:
			*data, err = b.ReadByte()
			if err != nil {
				return err
			}
		case *uint16:
			arr := b.Next(2)
			if len(arr) != 2 {
				return errors.New("too small")
			}
			*data = binary.LittleEndian.Uint16(arr)
		case *uint32:
			arr := b.Next(4)
			if len(arr) != 4 {
				return errors.New("too small")
			}
			*data = binary.LittleEndian.Uint32(arr)
		case *uint64:
			arr := b.Next(8)
			if len(arr) != 8 {
				return errors.New("too small")
			}
			*data = binary.LittleEndian.Uint64(arr)
		default:
			return errors.New("unknown type")
		}
	}
	return nil
}
