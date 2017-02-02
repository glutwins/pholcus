package master

import (
	"github.com/glutwins/pholcus/logs"
	"github.com/henrylee2cn/teleport"
)

// 主节点自动分配任务的操作
type masterTaskHandle struct {
}

func (self *masterTaskHandle) Process(receive *teleport.NetData) *teleport.NetData {
	return teleport.ReturnData(nil)
}

// 主节点自动接收从节点消息并打印的操作
type masterLogHandle struct{}

func (*masterLogHandle) Process(receive *teleport.NetData) *teleport.NetData {
	logs.Log.Informational(" * ")
	logs.Log.Informational(" *     [ %s ]    %s", receive.From, receive.Body)
	logs.Log.Informational(" * ")
	return nil
}
