package smart_client

import (
	"SmartHomeServer/define"
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
	itemType define.ItemType
	name     string
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
				fmt.Println(msg.Seq, msg.MsgId)
				if !this.ProcessMsg(msg.MsgId, msg.Data) {
					return
				}
				recvBuf.Next(int(needLen))
			} else {
				break
			}
		}
	}
}

func (this *SmartClient) ProcessMsg(msgId uint32, data []byte) bool {
	switch msgId {
	case uint32(message.MsgId_RegisterReq):
		this.handleRegister(data)
		break
	default:
		fmt.Println("没有找到消息:", msgId)
	}
	return true
}

func (this *SmartClient) Kick() {
	if this.registed {
		//todo
	}
	this.conn.Close()
}

func (this *SmartClient) sendAck(ack message.BaseAck) {
	data := ack.Encode()
	length := len(data)
	bf := &message.ByteBuffer{}
	bf.WriteInt32(int32(12 + length))
	bf.WriteInt32(0)
	bf.WriteInt32(ack.GetMsgId())
	bf.Write(data)
	this.conn.Write(bf.GetBuffer())
	fmt.Println(12+length, len(bf.GetBuffer()))
}

func (this *SmartClient) handleRegister(data []byte) {
	if this.registed {
		this.sendAck(&message.RegisterAck{Ret: message.Result_AlreaedRegistered})
		fmt.Println("设备注册过了")
		return
	}
	msg := &message.RegisterReq{}
	msg.Decode(data)
	this.itemType = define.ItemType(msg.ItemType)
	this.name = msg.Name
	this.registed = true
	this.sendAck(&message.RegisterAck{Ret: message.Result_Success})
	fmt.Println("设备注册成功")
}
