package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type BinaryReader struct {
	*bytes.Buffer
}

func (b *BinaryReader) ReadInt32() int32 {
	return int32(binary.LittleEndian.Uint32(b.Next(4)))
}
func (b *BinaryReader) ReadUnicode(length int) string {
	unicode := b.Next(2 * length)
	// 将读取的字节转换为字符串
	// 查找第一个null字符的索引
	index := bytes.IndexByte(unicode, '\x00')

	// 如果找到null字符,则返回null字符之前的子字符串
	if index > 0 {
		return string(unicode[:index])
	}
	return ""
}

func (b *BinaryReader) WriteUnicodeString(text string, length int) {
	// 将字符串转换为Unicode字节数组
	unicode := []byte(text)

	// 如果Unicode字节数组长度超过len*2,则截取前len*2个字节
	if length*2 < len(unicode) {
		unicode = unicode[:length*2]
	}

	// 写入Unicode字节数组到writer
	_, _ = b.Write(unicode)
}

func (b *BinaryReader) ReadToEnd() []byte {

	return b.Buffer.Next(b.Buffer.Len())
}

type MemoryStream struct {
	buff []byte
	loc  int
}

// DefaultCapacity is the size in bytes of a new MemoryStream's backing buffer
const DefaultCapacity = 512

// NewMemoryStream New creates a new MemoryStream instance
func NewMemoryStream() *MemoryStream {
	return NewCapacity(DefaultCapacity)
}

// NewCapacity starts the returned MemoryStream with the given capacity
func NewCapacity(cap int) *MemoryStream {
	return &MemoryStream{buff: make([]byte, 0, cap), loc: 0}
}

// Seek sets the offset for the next Read or Write to offset, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end. Seek
// returns the new offset and an error, if any.
//
// Seeking to a negative offset is an error. Seeking to any positive offset is
// legal. If the location is beyond the end of the current length, the position
// will be placed at length.
func (m *MemoryStream) Seek(offset int64, whence int) (int64, error) {
	newLoc := m.loc
	switch whence {
	case 0:
		newLoc = int(offset)
	case 1:
		newLoc += int(offset)
	case 2:
		newLoc = len(m.buff) - int(offset)
	}

	if newLoc < 0 {
		return int64(m.loc), errors.New("Unable to seek to a location <0")
	}

	if newLoc > len(m.buff) {
		newLoc = len(m.buff)
	}

	m.loc = newLoc

	return int64(m.loc), nil
}

// Read puts up to len(p) bytes into p. Will return the number of bytes read.
func (m *MemoryStream) Read(p []byte) (n int, err error) {
	n = copy(p, m.buff[m.loc:len(m.buff)])
	m.loc += n

	if m.loc == len(m.buff) {
		return n, io.EOF
	}

	return n, nil
}

// Write writes the given bytes into the memory stream. If needed, the underlying
// buffer will be expanded to fit the new bytes.
func (m *MemoryStream) Write(p []byte) (n int, err error) {
	// Do we have space?
	if available := cap(m.buff) - m.loc; available < len(p) {
		// How much should we expand by?
		addCap := cap(m.buff)
		if addCap < len(p) {
			addCap = len(p)
		}

		newBuff := make([]byte, len(m.buff), cap(m.buff)+addCap)

		copy(newBuff, m.buff)

		m.buff = newBuff
	}

	// Write
	n = copy(m.buff[m.loc:cap(m.buff)], p)
	m.loc += n
	if len(m.buff) < m.loc {
		m.buff = m.buff[:m.loc]
	}

	return n, nil
}

// Bytes returns a copy of ALL valid bytes in the stream, regardless of the current
// position.
func (m *MemoryStream) Bytes() []byte {
	b := make([]byte, len(m.buff))
	copy(b, m.buff)
	return b
}

// Rewind returns the stream to the beginning
func (m *MemoryStream) Rewind() (int64, error) {
	return m.Seek(0, 0)
}
