package slave

import (
	"encoding/json"

	"github.com/glutwins/pholcus/common/schema"
	"github.com/glutwins/pholcus/logs"
	"github.com/henrylee2cn/teleport"
)

type Slave struct {
	trans teleport.Teleport
}

func NewSlave(master string, port string) *Slave {
	m := &Slave{}
	m.trans = teleport.New()
	m.trans.SetAPI(teleport.API{
		"task": &slaveTaskHandle{},
	}).Client(master, port)

	go func() {
		for true {
			_, msg, ok := logs.Log.StealOne()
			if !ok {
				return
			}
			if self.Teleport.CountNodes() == 0 {
				// 与服务器失去连接后，抛掉返馈日志
				continue
			}
			self.Teleport.Request(msg, "log", "")
		}
	}()

	return m
}

// 从节点自动接收主节点任务的操作
type slaveTaskHandle struct {
}

func (self *slaveTaskHandle) Process(receive *teleport.NetData) *teleport.NetData {
	t := &schema.Task{}
	err := json.Unmarshal([]byte(receive.Body.(string)), t)
	if err != nil {
		logs.Log.Error("json解码失败 %v", receive.Body)
		return nil
	}
	return nil
}
