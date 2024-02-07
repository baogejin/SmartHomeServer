package message

import (
	"encoding/binary"
)

type ByteBuffer struct {
	pos int
	buf []byte
}

func (this *ByteBuffer) SetData(data []byte) {
	this.buf = data
	this.pos = 0
}

func (this *ByteBuffer) ReadInt32() int32 {
	if this.pos+4 > len(this.buf) {
		return 0
	}
	ret := binary.LittleEndian.Uint32(this.buf[this.pos:])
	this.pos += 4
	return int32(ret)
}

func (this *ByteBuffer) ReadString() string {
	if this.pos+4 > len(this.buf) {
		return ""
	}
	length := this.ReadInt32()
	if this.pos+int(length) > len(this.buf) {
		return ""
	}
	ret := string(this.buf[this.pos : this.pos+int(length)])
	this.pos += int(length)
	return ret
}

func (this *ByteBuffer) WriteInt32(v int32) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(v))
	this.buf = append(this.buf, data...)
}

func (this *ByteBuffer) WriteString(s string) {
	length := len(s)
	this.WriteInt32(int32(length))
	this.buf = append(this.buf, []byte(s)...)
}

func (this *ByteBuffer) Write(data []byte) {
	this.buf = append(this.buf, data...)
}

func (this *ByteBuffer) GetBuffer() []byte {
	return this.buf
}
