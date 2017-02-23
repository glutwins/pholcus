package slave

import (
	"encoding/json"
	"time"

	"github.com/glutwins/flow"
	"github.com/glutwins/pholcus/common/schema"
	"github.com/glutwins/pholcus/logs"
	"github.com/glutwins/pholcus/slave/scheduler"
	"github.com/henrylee2cn/teleport"
)

type CrawlTask struct {
	task   *schema.Task
	result *schema.TaskResult
	sdl    *scheduler.Scheduler
}

func (t *CrawlTask) Hash() int {
	return 0
}

func (t *CrawlTask) Exec() (interface{}, error) {
	sdl := scheduler.NewScheduler(t.task)
	t.result = &schema.TaskResult{
		StartTime: time.Now(),
	}

	go sdl.Run()

	return nil, nil
}

type Slave struct {
	trans teleport.Teleport
	tasks *flow.TaskFlow
}

func NewSlave(master string, port string, cnum int) *Slave {
	m := &Slave{}
	m.trans = teleport.New()
	m.trans.SetAPI(teleport.API{
		"task": &slaveTaskHandle{},
	}).Client(master, port)
	m.tasks = flow.NewTaskFlow(cnum)

	go func() {
		for true {
			_, msg, ok := logs.Log.StealOne()
			if !ok {
				return
			}
			if m.trans.CountNodes() == 0 {
				// 与服务器失去连接后，抛掉返馈日志
				continue
			}
			m.trans.Request(msg, "log", "")
		}
	}()

	return m
}

// 从节点自动接收主节点任务的操作
type slaveTaskHandle struct {
	s *Slave
}

func (self *slaveTaskHandle) Process(receive *teleport.NetData) *teleport.NetData {
	t := &schema.Task{}
	err := json.Unmarshal([]byte(receive.Body.(string)), t)
	if err != nil {
		logs.Log.Error("json解码失败 %v", receive.Body)
		return nil
	}

	self.s.tasks.Exec(&CrawlTask{task: t})
	return nil
}
