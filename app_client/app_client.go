package app_client

import (
	"SmartHomeServer/eventbus"
	"SmartHomeServer/message"
	"SmartHomeServer/smart_client"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type AppClient struct {
	conn net.Conn
	Uid  int32
}

func (this *AppClient) Start(conn net.Conn) {
	this.conn = conn
	buf := make([]byte, 2048)
	recvBuf := bytes.NewBuffer(make([]byte, 0, 2048))
	//注册事件
	eventbus.GetInstance().Subscribe(eventbus.Event_NewItem, this.onNewItem)
	defer eventbus.GetInstance().Unsubscribe(eventbus.Event_NewItem, this.onNewItem)

	eventbus.GetInstance().Subscribe(eventbus.Event_ItemChangeStatus, this.onItemChangeStatus)
	defer eventbus.GetInstance().Unsubscribe(eventbus.Event_NewItem, this.onItemChangeStatus)

	eventbus.GetInstance().Subscribe(eventbus.Event_ItemDisconnect, this.onItemDisconnect)
	defer eventbus.GetInstance().Unsubscribe(eventbus.Event_NewItem, this.onItemDisconnect)

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
				this.ProcessMsg(msg.MsgId, msg.Data)
				recvBuf.Next(int(needLen))
			} else {
				break
			}
		}
	}
}

func (this *AppClient) ProcessMsg(msgId uint32, data []byte) {
	switch msgId {
	case uint32(message.MsgId_GetSnartItemListReq):
		list := smart_client.GetClientMgr().GetSmartItemList()
		ack := &message.GetSmartItemListAck{
			Ret:   message.Result_Success,
			Items: list,
		}
		this.SendMsg(ack)
	case uint32(message.MsgId_SetSmartItemStatusReq):
		msg := &message.SetSmartItemStatusReq{}
		msg.Decode(data)
		ret := smart_client.GetClientMgr().SetSmartItemStatus(msg.Uid, msg.Status)
		ack := &message.SetSmartItemStatusAck{
			Ret: ret,
		}
		this.SendMsg(ack)
	default:
		fmt.Println("没有找到消息:", msgId)
	}
}

func (this *AppClient) SendMsg(msg message.BaseMsg) {
	data := msg.Encode()
	length := len(data)
	bf := &message.ByteBuffer{}
	bf.WriteInt32(int32(12 + length))
	bf.WriteInt32(0)
	bf.WriteInt32(msg.GetMsgId())
	bf.Write(data)
	this.conn.Write(bf.GetBuffer())
}

func (this *AppClient) onNewItem(item *message.SmartItem) {
	this.SendMsg(&message.NewItemNotify{
		Item: item,
	})
}

func (this *AppClient) onItemChangeStatus(uid, status int32) {
	this.SendMsg(&message.ItemChangeStatusNotify{
		Uid:    uid,
		Status: status,
	})
}

func (this *AppClient) onItemDisconnect(uid int32) {
	this.SendMsg(&message.ItemDisconnectNotify{
		Uid: uid,
	})
}
