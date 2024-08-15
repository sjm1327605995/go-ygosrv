package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"os"
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

// Buffer implements interfaces implemented by files.
// The main purpose of this type is to have an in memory replacement for a
// file.
type Buffer struct {
	// Buff is the backing buffer
	buff *bytes.Buffer
	// Index indicates where in the buffer we are at
	Index    int64
	isClosed bool
}

// New returns a new populated Buffer
func New(b []byte) *Buffer {
	return &Buffer{buff: bytes.NewBuffer(b)}
}

// NewFromReader is a convenience method that returns a new populated Buffer
// whose contents are sourced from a supplied reader by loading it entirely
// into memory.
func NewFromReader(reader io.Reader) (*Buffer, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return New(data), nil
}

// Bytes returns the bytes available until the end of the buffer.
func (f *Buffer) Bytes() []byte {
	if f.isClosed || f.Index >= int64(f.buff.Len()) {
		return []byte{}
	}
	return f.buff.Bytes()[f.Index:]
}

// String implements the Stringer interface
func (f *Buffer) String() string {
	return string(f.buff.Bytes()[f.Index:])
}

// Read implements io.Reader https://golang.org/pkg/io/#Reader
// Read reads up to len(p) bytes into p. It returns the number of bytes read (0 <= n <= len(p))
// and any error encountered. Even if Read returns n < len(p), it may use all of p as scratch
// space during the call. If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.

// When Read encounters an error or end-of-file condition after successfully reading n > 0 bytes,
// it returns the number of bytes read. It may return the (non-nil) error from the same call or
// return the error (and n == 0) from a subsequent call. An instance of this general case is
// that a Reader returning a non-zero number of bytes at the end of the input stream may return
// either err == EOF or err == nil. The next Read should return 0, EOF.
func (f *Buffer) Read(b []byte) (n int, err error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}
	if len(b) == 0 {
		return 0, nil
	}
	if f.Index >= int64(f.buff.Len()) {
		return 0, io.EOF
	}
	n, err = bytes.NewBuffer(f.buff.Bytes()[f.Index:]).Read(b)
	f.Index += int64(n)

	return n, err
}

// ReadAt implements io.ReaderAt https://golang.org/pkg/io/#ReaderAt
// ReadAt reads len(p) bytes into p starting at offset off in the underlying input source.
// It returns the number of bytes read (0 <= n <= len(p)) and any error encountered.
//
// When ReadAt returns n < len(p), it returns a non-nil error explaining why more bytes were not returned.
// In this respect, ReadAt is stricter than Read.
//
// Even if ReadAt returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, ReadAt blocks until either all the data is available or an error occurs.
// In this respect ReadAt is different from Read.
//
// If the n = len(p) bytes returned by ReadAt are at the end of the input source,
// ReadAt may return either err == EOF or err == nil.
//
// If ReadAt is reading from an input source with a seek offset,
// ReadAt should not affect nor be affected by the underlying seek offset.
// Clients of ReadAt can execute parallel ReadAt calls on the same input source.
func (f *Buffer) ReadAt(p []byte, off int64) (n int, err error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}
	if off < 0 {
		return 0, errors.New("filebuffer.ReadAt: negative offset")
	}
	reqLen := len(p)
	buffLen := int64(f.buff.Len())
	if off >= buffLen {
		return 0, io.EOF
	}

	n = copy(p, f.buff.Bytes()[off:])
	if n < reqLen {
		err = io.EOF
	}
	return n, err
}

// Write implements io.Writer https://golang.org/pkg/io/#Writer
// by appending the passed bytes to the buffer unless the buffer is closed or index negative.
func (f *Buffer) Write(p []byte) (n int, err error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}
	if f.Index < 0 {
		return 0, io.EOF
	}
	// we might have rewinded, let's reset the buffer before appending to it
	idx := int(f.Index)
	buffLen := f.buff.Len()
	if idx != buffLen && idx <= buffLen {
		f.buff = bytes.NewBuffer(f.Bytes()[:f.Index])
	}
	n, err = f.buff.Write(p)

	f.Index += int64(n)
	return n, err
}

// Seek implements io.Seeker https://golang.org/pkg/io/#Seeker
func (f *Buffer) Seek(offset int64, whence int) (idx int64, err error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}

	var abs int64
	switch whence {
	case 0:
		abs = offset
	case 1:
		abs = int64(f.Index) + offset
	case 2:
		abs = int64(f.buff.Len()) + offset
	default:
		return 0, errors.New("filebuffer.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("filebuffer.Seek: negative position")
	}
	f.Index = abs
	return abs, nil
}

// Close implements io.Closer https://golang.org/pkg/io/#Closer
// It closes the buffer, rendering it unusable for I/O. It returns an error, if any.
func (f *Buffer) Close() error {
	f.isClosed = true
	return nil
}
