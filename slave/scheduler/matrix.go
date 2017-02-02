package scheduler

import (
	"sort"
	"sync"
	"sync/atomic"

	"github.com/glutwins/pholcus/app/aid/history"
	"github.com/glutwins/pholcus/app/downloader/request"
	"github.com/glutwins/pholcus/app/spider"
	"github.com/glutwins/pholcus/common/schema"
	"github.com/glutwins/pholcus/logs"
)

// 一个Spider实例的请求矩阵
type Matrix struct {
	maxPage     int64                       // 最大采集页数，以负数形式表示
	sp          *spider.Spider              // 所属Spider
	task        *schema.Task                // 所属任务
	keyin       string                      // 自定义任务属性
	reqs        map[int][]*request.Request  // [优先级]队列，优先级默认为0
	priorities  []int                       // 优先级顺序，从低到高
	history     history.Historier           // 历史记录
	tempHistory map[string]bool             // 临时记录 [reqUnique(url+method)]true
	failures    map[string]*request.Request // 历史及本次失败请求

	sdl             *Scheduler // 所属调度器
	tempHistoryLock sync.RWMutex
	failureLock     sync.Mutex
	sync.Mutex
}

// 添加请求到队列，并发安全
func (self *Matrix) Push(req *request.Request) {
	// 禁止并发，降低请求积存量
	self.Lock()
	defer self.Unlock()

	// 达到请求上限，停止该规则运行
	if self.maxPage >= 0 {
		return
	}

	// 不可重复下载的req
	if !req.IsReloadable() {
		// 已存在成功记录时退出
		if self.hasHistory(req.Unique()) {
			return
		}
		// 添加到临时记录
		self.insertTempHistory(req.Unique())
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

	// 大致限制加入队列的请求量，并发情况下应该会比maxPage多
	atomic.AddInt64(&self.maxPage, 1)
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
			req.SetProxy(self.sdl.proxy.GetOne(req.GetUrl()))
			return req
		}
	}
	return nil
}

// 返回是否作为新的失败请求被添加至队列尾部
func (self *Matrix) DoHistory(req *request.Request, ok bool) bool {
	if !req.IsReloadable() {
		self.tempHistoryLock.Lock()
		delete(self.tempHistory, req.Unique())
		self.tempHistoryLock.Unlock()

		if ok {
			self.history.UpsertSuccess(req.Unique())
			return false
		}
	}

	if ok {
		return false
	}

	self.failureLock.Lock()
	defer self.failureLock.Unlock()
	if _, ok := self.failures[req.Unique()]; !ok {
		// 首次失败时，在任务队列末尾重新执行一次
		self.failures[req.Unique()] = req
		logs.Log.Informational(" *     + 失败请求: [%v]\n", req.GetUrl())
		return true
	}
	// 失败两次后，加入历史失败记录
	self.history.UpsertFailure(req)
	return false
}

// 非服务器模式下保存历史成功记录
func (self *Matrix) TryFlushSuccess() {
	self.history.FlushSuccess()
}

// 非服务器模式下保存历史失败记录
func (self *Matrix) TryFlushFailure() {
	self.history.FlushFailure()
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

func (self *Matrix) hasHistory(reqUnique string) bool {
	if self.history.HasSuccess(reqUnique) {
		return true
	}
	self.tempHistoryLock.RLock()
	has := self.tempHistory[reqUnique]
	self.tempHistoryLock.RUnlock()
	return has
}

func (self *Matrix) insertTempHistory(reqUnique string) {
	self.tempHistoryLock.Lock()
	self.tempHistory[reqUnique] = true
	self.tempHistoryLock.Unlock()
}

func (self *Matrix) setFailures(reqs map[string]*request.Request) {
	self.failureLock.Lock()
	defer self.failureLock.Unlock()
	for key, req := range reqs {
		self.failures[key] = req
		logs.Log.Informational(" *     + 失败请求: [%v]\n", req.GetUrl())
	}
}
