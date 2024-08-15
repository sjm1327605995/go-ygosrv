package utils

import (
	"bytes"
)

//func (cm *CoreMessage) SetEndPosition() {
//	cm.endPosition = 0 // This should be set to the actual end position in your context
//	cm.length = cm.endPosition - cm.startPosition
//}

type BinaryData interface {
	Raw() []byte
	//Marshal([]byte) error
	Unmarshal([]byte) error
}

type MessageData struct {
	raw []byte
	off int
}

func (m *MessageData) Raw() []byte {
	return m.raw
}
func (m *MessageData) ReadUnicode(length int) string {
	unicode := m.raw[m.off : 2*length]
	// 将读取的字节转换为字符串
	// 查找第一个null字符的索引
	index := bytes.IndexByte(unicode, '\x00')

	// 如果找到null字符,则返回null字符之前的子字符串
	if index > 0 {
		return string(unicode[:index])
	}
	return ""
}
