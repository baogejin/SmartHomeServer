package message

type MsgId int32

const (
	//设备消息id
	MsgId_RegisterReq      MsgId = 1
	MsgId_RegisterAck      MsgId = 2
	MsgId_ReportStatusReq  MsgId = 3
	MsgId_ReportStatusAck  MsgId = 4
	MsgId_ChangeStatusPush MsgId = 5

	//app消息id
	MsgId_GetSnartItemListReq    MsgId = 1001
	MsgId_GetSnartItemListAck    MsgId = 1002
	MsgId_SetSmartItemStatusReq  MsgId = 1003
	MsgId_SetSmartItemStatusAck  MsgId = 1004
	MsgId_ItemChangeStatusNotify MsgId = 1005
	MsgId_ItemDisconnectNotify   MsgId = 1006
	MsgId_NewItemNotify          MsgId = 1007
)
