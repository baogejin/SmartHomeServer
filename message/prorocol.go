package message

type BaseMsg interface {
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

type ReportStatusReq struct {
	Status int32
}

func (this *ReportStatusReq) Decode(data []byte) {
	bf := &ByteBuffer{}
	bf.SetData(data)
	this.Status = bf.ReadInt32()
}

type ReportStatusAck struct {
	Ret ResultId
}

func (this *ReportStatusAck) Encode() []byte {
	bf := &ByteBuffer{}
	bf.WriteInt32(int32(this.Ret))
	return bf.GetBuffer()
}

func (this *ReportStatusAck) GetMsgId() int32 {
	return int32(MsgId_ReportStatusAck)
}

type ChangeStatusPush struct {
	Status int32
}

func (this *ChangeStatusPush) Encode() []byte {
	bf := &ByteBuffer{}
	bf.WriteInt32(int32(this.Status))
	return bf.GetBuffer()
}

func (this *ChangeStatusPush) GetMsgId() int32 {
	return int32(MsgId_ChangeStatusPush)
}
