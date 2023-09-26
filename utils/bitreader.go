package utils

import (
	"encoding/binary"
	"errors"
)

type BitReader struct {
	arr []byte
	pos int
	len int
}

func NewBitReader(arr []byte, length int) *BitReader {
	return &BitReader{
		arr: arr,
		len: length,
	}
}
func (b *BitReader) ReadUint8() uint8 {
	b.pos += 1
	return b.arr[b.pos-1]
}
func (b *BitReader) ReadUint16() uint16 {
	val := binary.LittleEndian.Uint16(b.arr[b.pos:])
	b.pos += 2
	return val
}

func (b *BitReader) ReadUint32() uint32 {
	val := binary.LittleEndian.Uint32(b.arr[b.pos:])
	b.pos += 4
	return val
}
func (b *BitReader) ReadUint64() uint64 {
	val := binary.LittleEndian.Uint64(b.arr[b.pos:])
	b.pos += 8
	return val
}
func (b *BitReader) Next(size int) []byte {
	var (
		oldPos = b.pos
		total  = b.pos + size
		count  = len(b.arr)
	)
	if total > count {
		b.pos = count
	}
	defer func() {
		b.pos = 0
	}()
	return b.arr[oldPos:b.pos]
}
func (b *BitReader) PutOffsetUint8(offset int, val uint8) {
	b.arr[offset] = val
	b.pos += 1
}
func (b *BitReader) PutOffsetUint16(offset int, val uint16) {
	binary.LittleEndian.PutUint16(b.arr[offset:], val)
	b.pos += 2
}
func (b *BitReader) PutOffsetUint32(offset int, val uint32) {
	binary.LittleEndian.PutUint32(b.arr[offset:], val)
	b.pos += 4
}
func (b *BitReader) PutOffsetUint64(offset int, val uint64) {
	binary.LittleEndian.PutUint64(b.arr[offset:], val)
	b.pos += 8
}
func (b *BitReader) PutUint8(val uint8) {
	b.PutOffsetUint8(b.pos, val)
}
func (b *BitReader) PutUint16(val uint16) {
	b.PutOffsetUint16(b.pos, val)
}
func (b *BitReader) PutUint32(val uint32) {
	b.PutOffsetUint32(b.pos, val)
}
func (b *BitReader) PutUint64(val uint64) {
	b.PutOffsetUint64(b.pos, val)
}

func GetData(b *BitReader, ptrs ...interface{}) (err error) {

	for i := range ptrs {
		switch data := ptrs[i].(type) {
		case *uint8:
			*data = b.ReadUint8()
		case *uint16:
			*data = b.ReadUint16()
		case *uint32:
			*data = b.ReadUint32()
		case *uint64:
			*data = b.ReadUint64()
		default:
			return errors.New("unknown type")
		}
	}
	return nil
}
