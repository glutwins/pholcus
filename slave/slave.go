package slave

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sync"

	"github.com/glutwins/pholcus/common/schema"
	"github.com/glutwins/pholcus/logs"
	"github.com/henrylee2cn/teleport"
)

type Slave struct {
	base  int
	max   int
	trans teleport.Teleport
	tasks map[int]*TaskProcess
	sync.Mutex
	recall sync.Cond
}

func NewSlave(master string, port string, cnum int) *Slave {
	m := &Slave{
		max:   cnum,
		tasks: make(map[int]*TaskProcess),
		trans: teleport.New(),
	}

	m.trans.SetAPI(teleport.API{
		"task": &slaveTaskHandle{},
	}).Client(master, port)

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

func (s *Slave) RunTask(t *schema.Task) {
	s.Lock()
	n := len(s.tasks) + 1
	if n > math.MaxUint16 {
		panic("slave task overload")
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.base++
	tp := &TaskProcess{id: s.base, cancel: cancel}
	s.tasks[tp.id] = tp
	s.Unlock()

	go func() {
		if n > len(s.tasks) {
			s.recall.Wait()
		}
		result, err := tp.Run(ctx, t)
		fmt.Println(result, err)
		s.Lock()
		delete(s.tasks, tp.id)
		s.Unlock()
		s.recall.Signal()
	}()

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

	self.s.RunTask(t)
	return nil
}
