package message

import "encoding/binary"

type MsgPack struct {
	Seq   uint32
	MsgId uint32
	Data  []byte
}

func PackMsg(msgId uint32, data []byte) []byte {
	length := 4 + 8 + 4 + len(data)
	buf := make([]byte, length)
	binary.LittleEndian.PutUint32(buf, uint32(length))
	binary.LittleEndian.PutUint32(buf[4:], 0)
	binary.LittleEndian.PutUint32(buf[8:], msgId)
	if len(data) > 0 {
		copy(buf[12:], data)
	}
	return buf
}

func UnpackMsg(buf []byte) *MsgPack {
	msg := &MsgPack{}
	msg.Seq = binary.LittleEndian.Uint32(buf)
	msg.MsgId = binary.LittleEndian.Uint32(buf[4:])
	if len(buf) > 8 {
		msg.Data = buf[8:]
	}
	return msg
}
