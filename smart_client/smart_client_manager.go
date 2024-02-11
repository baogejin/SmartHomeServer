package smart_client

import (
	"SmartHomeServer/eventbus"
	"SmartHomeServer/message"
	"sync"
)

var clientMgr *SmartClientManager
var once sync.Once

type SmartClientManager struct {
	clients sync.Map
	uid     int32
}

func GetClientMgr() *SmartClientManager {
	once.Do(func() {
		clientMgr = &SmartClientManager{}
		clientMgr.init()
	})
	return clientMgr
}

func (this *SmartClientManager) init() {
	this.uid = 0
}

func (this *SmartClientManager) AddClient(client *SmartClient) {
	this.uid++
	client.Uid = this.uid
	this.clients.Store(this.uid, client)
	eventbus.GetInstance().Publish(eventbus.Event_NewItem, &message.SmartItem{
		Uid:    client.Uid,
		Name:   client.Name,
		Status: client.Status,
	})
}

func (this *SmartClientManager) DelClient(uid int32) {
	this.clients.Delete(uid)
	eventbus.GetInstance().Publish(eventbus.Event_ItemDisconnect, uid)
}

func (this *SmartClientManager) GetSmartItemList() []*message.SmartItem {
	ret := make([]*message.SmartItem, 0)
	this.clients.Range(func(key, value any) bool {
		uid := key.(int32)
		client := value.(*SmartClient)
		ret = append(ret, &message.SmartItem{
			Uid:    uid,
			Name:   client.Name,
			Status: client.Status,
		})
		return true
	})
	return ret
}

func (this *SmartClientManager) SetSmartItemStatus(uid, status int32) message.ResultId {
	if value, ok := this.clients.Load(uid); ok {
		client := value.(*SmartClient)
		client.ChangeStatus(status)
		return message.Result_Success
	}
	return message.Result_SmartItemNotFound
}
