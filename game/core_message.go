// This code is a translation from C# to Go, using standard libraries for binary operations and encoding.
package game

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"

	"io"
)

// BinaryExtensions provides methods to read and write Unicode strings.
type BinaryExtensions struct{}

// WriteUnicode writes a Unicode string to the provided writer.
func (be *BinaryExtensions) WriteUnicode(writer io.Writer, text string, length int) {
	unicode := utf16.Encode([]rune(text))
	result := make([]byte, length*2)
	max := length*2 - 2
	copy(result, uint16ToByte(unicode[:min(len(unicode), max/2)]))
	writer.Write(result)
}

// ReadUnicode reads a Unicode string from the provided reader.
func ReadUnicode(reader io.Reader, length int) (string, error) {
	unicode := make([]byte, length*2)
	_, err := reader.Read(unicode)
	if err != nil {
		return "", err
	}
	decoded := utf16.Decode(bytesToUint16(unicode))
	text := string(decoded)
	index := bytes.IndexRune([]byte(text), 0)
	if index != -1 {
		text = text[:index]
	}
	return text, nil
}

func uint16ToByte(input []uint16) []byte {
	buf := new(bytes.Buffer)
	for _, v := range input {
		binary.Write(buf, binary.LittleEndian, v)
	}
	return buf.Bytes()
}

func bytesToUint16(input []byte) []uint16 {
	result := make([]uint16, len(input)/2)
	for i := 0; i < len(input); i += 2 {
		result[i/2] = binary.LittleEndian.Uint16(input[i : i+2])
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
