package message

type BaseAck interface {
	Encode() []byte
	GetMsgId() int32
}

type RegisterReq struct {
	ItemType int32
	Name     string
}

func (this *RegisterReq) Decode(data []byte) {
	bf := &ByteBuffer{}
	bf.SetData(data)
	this.ItemType = bf.ReadInt32()
	this.Name = bf.ReadString()
}

type RegisterAck struct {
	Ret ResultId
}

func (this *RegisterAck) Encode() []byte {
	bf := &ByteBuffer{}
	bf.WriteInt32(int32(this.Ret))
	return bf.GetBuffer()
}

func (this *RegisterAck) GetMsgId() int32 {
	return int32(MsgId_RegisterAck)
}
