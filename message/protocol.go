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

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /////////////////////////以下是app的消息//////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type SmartItem struct {
	Uid    int32
	Name   string
	Status int32
}

func (this *SmartItem) Encode() []byte {
	bf := &ByteBuffer{}
	bf.WriteInt32(this.Uid)
	bf.WriteString(this.Name)
	bf.WriteInt32(this.Status)
	return bf.GetBuffer()
}

type GetSmartItemListReq struct {
}

func (this *GetSmartItemListReq) Decode(data []byte) {
	//do nothing
}

type GetSmartItemListAck struct {
	Ret   ResultId
	Items []*SmartItem
}

func (this *GetSmartItemListAck) Encode() []byte {
	bf := &ByteBuffer{}
	bf.WriteInt32(int32(this.Ret))
	bf.WriteInt32(int32(len(this.Items)))
	for _, item := range this.Items {
		bf.Write(item.Encode())
	}
	return bf.GetBuffer()
}

func (this *GetSmartItemListAck) GetMsgId() int32 {
	return int32(MsgId_GetSmartItemListAck)
}

type SetSmartItemStatusReq struct {
	Uid    int32
	Status int32
}

func (this *SetSmartItemStatusReq) Decode(data []byte) {
	bf := &ByteBuffer{}
	bf.SetData(data)
	this.Uid = bf.ReadInt32()
	this.Status = bf.ReadInt32()
}

type SetSmartItemStatusAck struct {
	Ret ResultId
}

func (this *SetSmartItemStatusAck) Encode() []byte {
	bf := &ByteBuffer{}
	bf.WriteInt32(int32(this.Ret))
	return bf.GetBuffer()
}

func (this *SetSmartItemStatusAck) GetMsgId() int32 {
	return int32(MsgId_SetSmartItemStatusAck)
}

type ItemChangeStatusNotify struct {
	Uid    int32
	Status int32
}

func (this *ItemChangeStatusNotify) Encode() []byte {
	bf := &ByteBuffer{}
	bf.WriteInt32(this.Uid)
	bf.WriteInt32(this.Status)
	return bf.GetBuffer()
}

func (this *ItemChangeStatusNotify) GetMsgId() int32 {
	return int32(MsgId_ItemChangeStatusNotify)
}

type ItemDisconnectNotify struct {
	Uid int32
}

func (this *ItemDisconnectNotify) Encode() []byte {
	bf := &ByteBuffer{}
	bf.WriteInt32(this.Uid)
	return bf.GetBuffer()
}

func (this *ItemDisconnectNotify) GetMsgId() int32 {
	return int32(MsgId_ItemDisconnectNotify)
}

type NewItemNotify struct {
	Item *SmartItem
}

func (this *NewItemNotify) Encode() []byte {
	bf := &ByteBuffer{}
	bf.Write(this.Item.Encode())
	return bf.GetBuffer()
}

func (this *NewItemNotify) GetMsgId() int32 {
	return int32(MsgId_NewItemNotify)
}
