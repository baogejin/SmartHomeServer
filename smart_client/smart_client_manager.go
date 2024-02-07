package smart_client

import "sync"

var clientMgr *SmartClientManager
var once sync.Once

type SmartClientManager struct {
	clients sync.Map
	uid     uint64
}

func Get() *SmartClientManager {
	once.Do(func() {
		clientMgr = &SmartClientManager{}
		clientMgr.init()
	})
	return clientMgr
}

func (this *SmartClientManager) init() {
	this.uid = 0
}

func (this *SmartClientManager) AddClient() {

}
