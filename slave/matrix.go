package slave

import (
	"context"
	"sort"
	"sync"

	"github.com/glutwins/pholcus/app/aid/history"
	"github.com/glutwins/pholcus/app/downloader"
	"github.com/glutwins/pholcus/app/downloader/request"
	"github.com/glutwins/pholcus/app/pipeline"
	"github.com/glutwins/pholcus/app/spider"
	"github.com/glutwins/pholcus/common/schema"
)

// 一个Spider实例的请求矩阵
type Matrix struct {
	sp         *spider.Spider             // 所属Spider
	task       *schema.Task               // 所属任务
	reqs       map[int][]*request.Request // [优先级]队列，优先级默认为0
	priorities []int                      // 优先级顺序，从低到高
	history    *history.History           // 历史记录
	pipe       pipeline.Pipeline

	sync.Mutex
}

func NewMatrix(name string, keyin string, task *schema.Task, sp *spider.Spider) *Matrix {
	m := &Matrix{
		sp:         sp,
		reqs:       make(map[int][]*request.Request),
		priorities: []int{},
		history:    history.New(sp.GetName(), sp.GetSubName(), task.Db),
		pipe:       pipeline.New(sp, task),
	}
	return m
}

func (m *Matrix) Run(ctx context.Context) {
	for req := m.Pull(); req != nil; {
		if resp, err := downloader.Download(req); err != nil {
			m.history.Upsert(req, err)
		} else {
			ctx := spider.GetContext(m.sp, req)
			ctx.Response = resp
			// 过程处理，提炼数据
			ctx.Parse(req.GetRuleName())

			// 该条请求文件结果存入pipeline
			for _, f := range ctx.PullFiles() {
				if m.pipe.CollectFile(f) != nil {
					break
				}
			}
			// 该条请求文本结果存入pipeline
			for _, item := range ctx.PullItems() {
				if m.pipe.CollectData(item) != nil {
					break
				}
			}

			m.history.Upsert(req, nil)

			// 释放ctx准备复用
			spider.PutContext(ctx)
		}
	}
}

// 添加请求到队列，并发安全
func (self *Matrix) Push(req *request.Request) {
	// 禁止并发，降低请求积存量
	self.Lock()
	defer self.Unlock()

	if self.history.HasSucc(req.Unique()) {
		return
	}

	var priority = req.GetPriority()

	// 初始化该蜘蛛下该优先级队列
	if _, found := self.reqs[priority]; !found {
		self.priorities = append(self.priorities, priority)
		sort.Ints(self.priorities) // 从小到大排序
		self.reqs[priority] = []*request.Request{}
	}

	// 添加请求到队列
	self.reqs[priority] = append(self.reqs[priority], req)
}

// 从队列取出请求，不存在时返回nil，并发安全
func (self *Matrix) Pull() *request.Request {
	self.Lock()
	defer self.Unlock()
	// 按优先级从高到低取出请求
	for i := len(self.reqs) - 1; i >= 0; i-- {
		idx := self.priorities[i]
		if len(self.reqs[idx]) > 0 {
			req := self.reqs[idx][0]
			self.reqs[idx] = self.reqs[idx][1:]
			return req
		}
	}
	return nil
}

func (self *Matrix) TryFlush() {
	self.history.Flush()
}

func (self *Matrix) Len() int {
	self.Lock()
	defer self.Unlock()
	var l int
	for _, reqs := range self.reqs {
		l += len(reqs)
	}
	return l
}
