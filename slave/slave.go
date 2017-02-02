package slave

import (
	"encoding/json"
	"time"

	"github.com/glutwins/flow"
	"github.com/glutwins/pholcus/common/schema"
	"github.com/glutwins/pholcus/logs"
	"github.com/glutwins/pholcus/slave/matrix"
	"github.com/henrylee2cn/teleport"
)

type CrawlTask struct {
	task   *schema.Task
	result *schema.TaskResult
}

func (t *CrawlTask) Hash() int {
	return 0
}

func (t *CrawlTask) Exec() (interface{}, error) {
	matrix.NewMatrix(t.task)
	t.result = &schema.TaskResult{
		StartTime: time.Now(),
	}

	// 初始化资源队列
	scheduler.Init()

	// 执行任务
	var i int
	for i = 0; i < count && self.status != status.STOP; i++ {
	pause:
		if self.status == status.PAUSE {
			time.Sleep(time.Second)
			goto pause
		}
		// 从爬行队列取出空闲蜘蛛，并发执行
		c := self.CrawlerPool.Use()
		if c != nil {
			go func(i int, c crawler.Crawler) {
				// 执行并返回结果消息
				c.Init(sq.GetByIndex(i)).Run()
				// 任务结束后回收该蜘蛛
				self.RWMutex.RLock()
				if self.status != status.STOP {
					self.CrawlerPool.Free(c)
				}
				self.RWMutex.RUnlock()
			}(i, c)
		}
	}
	// 监控结束任务
	for ii := 0; ii < i; ii++ {
		s := <-cache.ReportChan
		if (s.DataNum == 0) && (s.FileNum == 0) {
			logs.Log.Informational(" *     [任务小计：%s | KEYIN：%s]   无采集结果，用时 %v！\n", s.SpiderName, s.Keyin, s.Time)
			continue
		}

		self.sum[0] += s.DataNum
		self.sum[1] += s.FileNum
	}
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

	self.s.tasks.Exec(&CrawlTask{t})
	return nil
}
