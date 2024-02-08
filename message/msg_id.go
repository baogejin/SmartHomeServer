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
)
