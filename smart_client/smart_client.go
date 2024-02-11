package smart_client

import (
	"SmartHomeServer/define"
	"SmartHomeServer/eventbus"
	"SmartHomeServer/message"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type SmartClient struct {
	conn     net.Conn
	seq      uint32
	registed bool
	ItemType define.ItemType
	Name     string
	Uid      int32
	Status   int32
}

func (this *SmartClient) Start(c net.Conn) {
	this.conn = c
	this.seq = 0
	this.registed = false
	buf := make([]byte, 2048)
	recvBuf := bytes.NewBuffer(make([]byte, 0, 2048))
	this.seq = 0
	for {
		length, err := this.conn.Read(buf)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("read err:", err)
			return
		}
		if length == 0 {
			fmt.Println("消息超过了缓冲长度")
			return
		}
		recvBuf.Write(buf[:length])
		for recvBuf.Len() > 4 {
			needLen := binary.LittleEndian.Uint32(recvBuf.Bytes())
			if recvBuf.Len() >= int(needLen) {
				msg := message.UnpackMsg(recvBuf.Bytes()[4:needLen])
				this.seq++
				if this.seq > define.MAX_SEQ {
					this.seq = this.seq % define.MAX_SEQ
				}
				if this.seq != msg.Seq {
					//序列不对
					this.Kick()
					break
				}
				if !this.registed && msg.MsgId != uint32(message.MsgId_RegisterReq) {
					//设备的第一条消息必须是注册
					this.Kick()
					break
				}
				this.ProcessMsg(msg.MsgId, msg.Data)
				recvBuf.Next(int(needLen))
			} else {
				break
			}
		}
	}
}

func (this *SmartClient) ProcessMsg(msgId uint32, data []byte) {
	switch msgId {
	case uint32(message.MsgId_RegisterReq):
		this.handleRegister(data)
	case uint32(message.MsgId_ReportStatusReq):
		this.handleReportStatus(data)
	default:
		fmt.Println("没有找到消息:", msgId)
	}
}

func (this *SmartClient) Kick() {
	if this.registed && this.Uid > 0 {
		GetClientMgr().DelClient(this.Uid)
	}
	this.conn.Close()
}

func (this *SmartClient) sendMsg(msg message.BaseMsg) {
	data := msg.Encode()
	length := len(data)
	bf := &message.ByteBuffer{}
	bf.WriteInt32(int32(12 + length))
	bf.WriteInt32(0)
	bf.WriteInt32(msg.GetMsgId())
	bf.Write(data)
	this.conn.Write(bf.GetBuffer())
}

func (this *SmartClient) handleRegister(data []byte) {
	if this.registed {
		this.sendMsg(&message.RegisterAck{Ret: message.Result_AlreaedRegistered})
		fmt.Println("设备注册过了")
		return
	}
	msg := &message.RegisterReq{}
	msg.Decode(data)
	this.ItemType = define.ItemType(msg.ItemType)
	this.Name = msg.Name
	this.registed = true
	this.sendMsg(&message.RegisterAck{Ret: message.Result_Success})
	GetClientMgr().AddClient(this)
	fmt.Println("设备注册成功")
}

func (this *SmartClient) handleReportStatus(data []byte) {
	msg := &message.ReportStatusReq{}
	msg.Decode(data)
	fmt.Println("状态上报，名称:", this.Name, ",状态:", msg.Status)
	this.Status = msg.Status
	eventbus.GetInstance().Publish(eventbus.Event_ItemChangeStatus, this.Uid, this.Status)
	this.sendMsg(&message.ReportStatusAck{Ret: message.Result_Success})
}

func (this *SmartClient) ChangeStatus(status int32) {
	fmt.Println("正在改变物品状态，名称:", this.Name, "目标状态:", status)
	msg := &message.ChangeStatusPush{Status: status}
	this.sendMsg(msg)
}
