// app interface for graphical user interface.
// 基本业务执行顺序依次为：New()-->[SetLog(io.Writer)-->]Init()-->SpiderPrepare()-->Run()
package app

import (
	"sync"
	"time"

	"github.com/glutwins/pholcus/app/crawler"
	"github.com/glutwins/pholcus/app/distribute"
	"github.com/glutwins/pholcus/app/scheduler"
	"github.com/glutwins/pholcus/app/spider"
	"github.com/glutwins/pholcus/logs"
	"github.com/glutwins/pholcus/runtime/cache"
	"github.com/glutwins/pholcus/runtime/status"
	"github.com/henrylee2cn/teleport"
)

type (
	App interface {
		Init() App                                   // 使用App前必须进行先Init初始化，SetLog()除外
		SpiderPrepare(original []*spider.Spider) App // 须在设置全局运行参数后Run()前调用（client模式下不调用该方法）
		Run()                                        // 阻塞式运行直至任务完成（须在所有应当配置项配置完成后调用）
		Stop()                                       // Offline 模式下中途终止任务（对外为阻塞式运行直至当前任务终止）
		IsRunning() bool                             // 检查任务是否正在运行
		IsPause() bool                               // 检查任务是否处于暂停状态
		IsStopped() bool                             // 检查任务是否已经终止
		PauseRecover()                               // Offline 模式下暂停\恢复任务
		Status() int                                 // 返回当前状态
		GetSpiderLib() []*spider.Spider              // 获取全部蜘蛛种类
		GetSpiderByName(string) *spider.Spider       // 通过名字获取某蜘蛛
		GetSpiderQueue() crawler.SpiderQueue         // 获取蜘蛛队列接口实例
	}
	Logic struct {
		*cache.AppConf                      // 全局配置
		*spider.SpiderSpecies               // 全部蜘蛛种类
		crawler.SpiderQueue                 // 当前任务的蜘蛛队列
		*distribute.TaskJar                 // 服务器与客户端间传递任务的存储库
		crawler.CrawlerPool                 // 爬行回收池
		teleport.Teleport                   // socket长连接双工通信接口，json数据传输
		sum                   [2]uint64     // 执行计数
		takeTime              time.Duration // 执行计时
		status                int           // 运行状态
		canSocketLog          bool
		sync.RWMutex
	}
)

// 全局唯一的核心接口实例
var LogicApp = New()

func New() App {
	return &Logic{
		AppConf:       cache.Task,
		SpiderSpecies: spider.Species,
		status:        status.STOPPED,
		Teleport:      teleport.New(),
		TaskJar:       distribute.NewTaskJar(),
		SpiderQueue:   crawler.NewSpiderQueue(),
		CrawlerPool:   crawler.NewCrawlerPool(),
	}
}

// 使用App前必须先进行Init初始化（SetLog()除外）
func (self *Logic) Init() App {
	self.Teleport = teleport.New()
	self.TaskJar = distribute.NewTaskJar()
	self.SpiderQueue = crawler.NewSpiderQueue()
	self.CrawlerPool = crawler.NewCrawlerPool()
	return self
}

// SpiderPrepare()必须在设置全局运行参数之后，Run()的前一刻执行
// original为spider包中未有过赋值操作的原始蜘蛛种类
// 已被显式赋值过的spider将不再重新分配Keyin
// client模式下不调用该方法
func (self *Logic) SpiderPrepare(original []*spider.Spider) App {
	self.SpiderQueue.Reset()
	// 遍历任务
	for _, sp := range original {
		spcopy := sp.Copy()
		spcopy.SetPausetime(self.AppConf.Pausetime)
		if spcopy.GetLimit() == spider.LIMIT {
			spcopy.SetLimit(self.AppConf.Limit)
		} else {
			spcopy.SetLimit(-1 * self.AppConf.Limit)
		}
		self.SpiderQueue.Add(spcopy)
	}
	// 遍历自定义配置
	self.SpiderQueue.AddKeyins(self.AppConf.Keyins)
	return self
}

// 获取全部蜘蛛种类
func (self *Logic) GetSpiderLib() []*spider.Spider {
	return self.SpiderSpecies.Get()
}

// 通过名字获取某蜘蛛
func (self *Logic) GetSpiderByName(name string) *spider.Spider {
	return self.SpiderSpecies.GetByName(name)
}

// 服务器客户端模式下返回节点数
func (self *Logic) CountNodes() int {
	return self.Teleport.CountNodes()
}

// 获取蜘蛛队列接口实例
func (self *Logic) GetSpiderQueue() crawler.SpiderQueue {
	return self.SpiderQueue
}

// 运行任务
func (self *Logic) Run() {
	// 确保开启报告
	logs.Log.GoOn()
	// 重置计数
	self.sum[0], self.sum[1] = 0, 0
	// 重置计时
	self.takeTime = 0
	self.client()
}

// Offline 模式下暂停\恢复任务
func (self *Logic) PauseRecover() {
	switch self.Status() {
	case status.PAUSE:
		self.setStatus(status.RUN)
	case status.RUN:
		self.setStatus(status.PAUSE)
	}

	scheduler.PauseRecover()
}

// Offline 模式下中途终止任务
func (self *Logic) Stop() {
	if self.status == status.STOPPED {
		return
	}
	if self.status != status.STOP {
		// 不可颠倒停止的顺序
		self.setStatus(status.STOP)
		scheduler.Stop()
		self.CrawlerPool.Stop()
	}
	for !self.IsStopped() {
		time.Sleep(time.Second)
	}
}

// 检查任务是否正在运行
func (self *Logic) IsRunning() bool {
	return self.status == status.RUN
}

// 检查任务是否处于暂停状态
func (self *Logic) IsPause() bool {
	return self.status == status.PAUSE
}

// 检查任务是否已经终止
func (self *Logic) IsStopped() bool {
	return self.status == status.STOPPED
}

// 返回当前运行状态
func (self *Logic) Status() int {
	self.RWMutex.RLock()
	defer self.RWMutex.RUnlock()
	return self.status
}

// 返回当前运行状态
func (self *Logic) setStatus(status int) {
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()
	self.status = status
}

// 客户端模式运行
func (self *Logic) client() {
	for {
		// 从任务库获取一个任务
		t := self.TaskJar.Pull()

		if self.Status() == status.STOP || self.Status() == status.STOPPED {
			return
		}

		self.SpiderQueue.Reset()

		self.AppConf.ThreadNum = t.ThreadNum
		self.AppConf.Pausetime = t.Pausetime
		self.AppConf.DockerCap = t.DockerCap
		self.AppConf.Inherit = t.Inherit
		self.AppConf.Limit = t.Limit
		self.AppConf.ProxyMinute = t.ProxyMinute
		self.AppConf.Keyins = t.Keyins

		// 初始化蜘蛛队列
		for _, n := range t.Spiders {
			sp := self.GetSpiderByName(n["name"])
			if sp == nil {
				continue
			}
			spcopy := sp.Copy()
			spcopy.SetPausetime(t.Pausetime)
			if spcopy.GetLimit() > 0 {
				spcopy.SetLimit(t.Limit)
			} else {
				spcopy.SetLimit(-1 * t.Limit)
			}
			if v, ok := n["keyin"]; ok {
				spcopy.SetKeyin(v)
			}
			self.SpiderQueue.Add(spcopy)
		}

		// 重置计数
		self.sum[0], self.sum[1] = 0, 0
		// 重置计时
		self.takeTime = 0

		count := self.SpiderQueue.Len()
		cache.ResetPageCount()
		// 初始化资源队列
		scheduler.Init()

		// 设置爬虫队列
		crawlerCap := self.CrawlerPool.Reset(count)

		logs.Log.Informational(" *     执行任务总数(任务数[*自定义配置数])为 %v 个\n", count)
		logs.Log.Informational(" *     采集引擎池容量为 %v\n", crawlerCap)
		logs.Log.Informational(" *     并发协程最多 %v 个\n", self.AppConf.ThreadNum)
		logs.Log.Informational(" *     随机停顿区间为 %v~%v 毫秒\n", self.AppConf.Pausetime/2, self.AppConf.Pausetime*2)
		logs.Log.Informational(" *                                                                                                 —— 开始抓取，请耐心等候 ——")
		logs.Log.Informational(` *********************************************************************************************************************************** `)

		// 开始计时
		cache.StartTime = time.Now()

		// 执行任务
		var i int
		for i = 0; i < count && self.Status() != status.STOP; i++ {
		pause:
			if self.IsPause() {
				time.Sleep(time.Second)
				goto pause
			}
			// 从爬行队列取出空闲蜘蛛，并发执行
			c := self.CrawlerPool.Use()
			if c != nil {
				go func(i int, c crawler.Crawler) {
					// 执行并返回结果消息
					c.Init(self.SpiderQueue.GetByIndex(i)).Run()
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
			logs.Log.Informational(" * ")
			switch {
			case s.DataNum > 0 && s.FileNum == 0:
				logs.Log.Informational(" *     [任务小计：%s | KEYIN：%s]   共采集数据 %v 条，用时 %v！\n",
					s.SpiderName, s.Keyin, s.DataNum, s.Time)
			case s.DataNum == 0 && s.FileNum > 0:
				logs.Log.Informational(" *     [任务小计：%s | KEYIN：%s]   共下载文件 %v 个，用时 %v！\n",
					s.SpiderName, s.Keyin, s.FileNum, s.Time)
			default:
				logs.Log.Informational(" *     [任务小计：%s | KEYIN：%s]   共采集数据 %v 条 + 下载文件 %v 个，用时 %v！\n",
					s.SpiderName, s.Keyin, s.DataNum, s.FileNum, s.Time)
			}

			self.sum[0] += s.DataNum
			self.sum[1] += s.FileNum
		}

		// 总耗时
		self.takeTime = time.Since(cache.StartTime)
		var prefix = func() string {
			if self.Status() == status.STOP {
				return "任务中途取消："
			}
			return "本次"
		}()
		// 打印总结报告
		logs.Log.Informational(" * ")
		logs.Log.Informational(` *********************************************************************************************************************************** `)
		logs.Log.Informational(" * ")
		switch {
		case self.sum[0] > 0 && self.sum[1] == 0:
			logs.Log.Informational(" *                            —— %s合计采集【数据 %v 条】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
				prefix, self.sum[0], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), self.takeTime)
		case self.sum[0] == 0 && self.sum[1] > 0:
			logs.Log.Informational(" *                            —— %s合计采集【文件 %v 个】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
				prefix, self.sum[1], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), self.takeTime)
		case self.sum[0] == 0 && self.sum[1] == 0:
			logs.Log.Informational(" *                            —— %s无采集结果，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
				prefix, cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), self.takeTime)
		default:
			logs.Log.Informational(" *                            —— %s合计采集【数据 %v 条 + 文件 %v 个】，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
				prefix, self.sum[0], self.sum[1], cache.GetPageCount(1), cache.GetPageCount(-1), cache.GetPageCount(0), self.takeTime)
		}
		logs.Log.Informational(" * ")
		logs.Log.Informational(` *********************************************************************************************************************************** `)

	}
}
