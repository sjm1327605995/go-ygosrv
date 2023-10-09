package websocket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"go-ygosrv/core/duel"
	"io"
)

type WsDecoder struct {
	upgraded bool
	buf      *bytes.Buffer
	wsMsgBuf wsMessageBuf // ws 消息缓存

}
type wsMessageBuf struct {
	firstHeader *ws.Header
	curHeader   *ws.Header
	cachedBuf   bytes.Buffer
}
type readWrite struct {
	io.Reader
	io.Writer
}

var (
	DataLoseError = errors.New("lost data")
)

func (w *WsDecoder) Decode(buff *bytes.Buffer, player *duel.DuelPlayer) gnet.Action {
	ok := w.upgrade(player.Conn, buff)
	if !ok {
		return gnet.Close
	}
	if buff.Len() <= 0 {
		return gnet.None
	}
	messages, err := w.readWsMessages()
	if err != nil {
		logging.Infof("Error reading message! %v", err)
		return gnet.None
	}
	if messages == nil || len(messages) <= 0 { //没有读到完整数据 不处理
		return gnet.None
	}
	for _, message := range messages {
		if message.OpCode.IsControl() {
			err = wsutil.HandleClientControlMessage(player.Conn, message)
			if err != nil {
				return gnet.None
			}
			continue
		}
		if message.OpCode == ws.OpBinary {

			duel.HandleCTOSPacket(player, message.Payload[2:], binary.LittleEndian.Uint16(message.Payload))
		}
	}
	return gnet.None
}
func (w *WsDecoder) upgrade(c gnet.Conn, buf *bytes.Buffer) (ok bool) {
	if w.upgraded {
		ok = true
		return
	}
	w.buf = buf
	tmpReader := bytes.NewReader(buf.Bytes())
	oldLen := tmpReader.Len()
	logging.Infof("do Upgrade")

	hs, err := ws.Upgrade(readWrite{tmpReader, c})
	skipN := oldLen - tmpReader.Len()
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF { //数据不完整
			return
		}
		buf.Next(skipN)
		logging.Infof("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
		return
	}
	buf.Next(skipN)
	logging.Infof("conn[%v] upgrade websocket protocol! Handshake: %v", c.RemoteAddr().String(), hs)
	if err != nil {
		logging.Infof("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
		return
	}
	ok = true
	w.upgraded = true
	return
}
func (w *WsDecoder) readBufferBytes(c gnet.Conn) gnet.Action {
	size := c.InboundBuffered()
	buf := make([]byte, size, size)
	read, err := c.Read(buf)
	if err != nil {
		logging.Infof("read err! %w", err)
		return gnet.Close
	}
	if read < size {
		logging.Infof("read bytes len err! size: %d read: %d", size, read)
		return gnet.Close
	}
	w.buf.Write(buf)
	return gnet.None
}

func (w *WsDecoder) readWsMessages() (messages []wsutil.Message, err error) {
	msgBuf := &w.wsMsgBuf
	in := w.buf
	for {
		if msgBuf.curHeader == nil {
			if in.Len() < ws.MinHeaderSize { //头长度至少是2
				return
			}
			var head ws.Header
			if in.Len() >= ws.MaxHeaderSize {
				head, err = ws.ReadHeader(in)
				if err != nil {
					return messages, err
				}
			} else { //有可能不完整，构建新的 reader 读取 head 读取成功才实际对 in 进行读操作
				tmpReader := bytes.NewReader(in.Bytes())
				oldLen := tmpReader.Len()
				head, err = ws.ReadHeader(tmpReader)
				skipN := oldLen - tmpReader.Len()
				if err != nil {
					if err == io.EOF || err == io.ErrUnexpectedEOF { //数据不完整
						return messages, nil
					}
					in.Next(skipN)
					return nil, err
				}
				in.Next(skipN)
			}

			msgBuf.curHeader = &head
			err = ws.WriteHeader(&msgBuf.cachedBuf, head)
			if err != nil {
				return nil, err
			}
		}
		dataLen := (int)(msgBuf.curHeader.Length)
		if dataLen > 0 {
			if in.Len() >= dataLen {
				_, err = io.CopyN(&msgBuf.cachedBuf, in, int64(dataLen))
				if err != nil {
					return
				}
			} else { //数据不完整
				fmt.Println(in.Len(), dataLen)
				logging.Infof("incomplete data")
				return
			}
		}
		if msgBuf.curHeader.Fin { //当前 header 已经是一个完整消息
			messages, err = wsutil.ReadClientMessage(&msgBuf.cachedBuf, messages)
			if err != nil {
				return nil, err
			}
			msgBuf.cachedBuf.Reset()
		} else {
			logging.Infof("The data is split into multiple frames")
		}
		msgBuf.curHeader = nil
	}
}
