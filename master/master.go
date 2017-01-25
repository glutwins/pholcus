package master

import (
	"github.com/glutwins/pholcus/common/schema"
	"github.com/henrylee2cn/teleport"
)

type Master struct {
	trans teleport.Teleport
}

func NewMaster(addr string) *Master {
	m := &Master{}
	m.trans = teleport.New()
	m.trans.SetAPI(teleport.API{
		"task": &masterTaskHandle{},
		"log":  &masterLogHandle{},
	}).Server(addr)

	return m
}

func (m *Master) AddTask(t *schema.Task) {

	self.TaskJar.Push(&t)
}

func (m *Master) ListTasks() {

}

func (m *Master) ShowTaskLog() {

}
