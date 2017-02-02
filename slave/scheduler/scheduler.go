package scheduler

import (
	"sync/atomic"
	"time"

	"github.com/glutwins/flow"
	"github.com/glutwins/pholcus/app/aid/history"
	"github.com/glutwins/pholcus/app/aid/proxy"
	"github.com/glutwins/pholcus/app/downloader"
	"github.com/glutwins/pholcus/app/downloader/request"
	"github.com/glutwins/pholcus/app/pipeline"
	"github.com/glutwins/pholcus/app/spider"
	"github.com/glutwins/pholcus/common/schema"
	"github.com/glutwins/pholcus/runtime/status"
)

// 调度器
type Scheduler struct {
	task     *schema.Task
	status   int32        // 运行状态
	proxy    *proxy.Proxy // 全局代理IP
	matrices []*Matrix    // Spider实例的请求矩阵列表
	threads  *flow.TaskFlow

	down downloader.Downloader        //全局公用的下载器
	pipe map[string]pipeline.Pipeline //结果收集与输出管道
}

func NewScheduler(task *schema.Task) *Scheduler {
	sdl := &Scheduler{
		proxy:   proxy.New(),
		status:  status.RUN,
		threads: flow.NewTaskFlow(task.ThreadNum),
		pipe:    map[string]pipeline.Pipeline{},
		down:    downloader.SurferDownloader,
	}
	for name, keyins := range task.Spiders {
		sp := spider.Species.GetByName(name)
		sdl.pipe[name] = pipeline.New(sp, task)
		for _, keyin := range keyins {
			m := &Matrix{
				sp:          sp,
				maxPage:     task.Limit,
				reqs:        make(map[int][]*request.Request),
				priorities:  []int{},
				history:     history.New(sp.GetName(), sp.GetSubName(), task.Db),
				tempHistory: make(map[string]bool),
				failures:    make(map[string]*request.Request),
				keyin:       keyin,
			}

			m.history.ReadSuccess()
			m.history.ReadFailure()
			m.setFailures(m.history.PullFailure())
			sdl.matrices = append(sdl.matrices, m)
		}
	}
	sdl.proxy.UpdateTicker(task.ProxyMinute)
	return sdl
}

func (self *Scheduler) Run() {
	for _, pipe := range self.pipe {
		pipe.Start()
	}

	for name, _ := range self.task.Spiders {
		if sp := spider.Species.GetByName(name); sp != nil {
			sp.RuleTree.Root(spider.GetContext(sp, nil))
		}
	}

	for {
		st := atomic.LoadInt32(&self.status)
		if st == status.STOP {
			for self.threads.JobCount() > 0 {
				time.Sleep(time.Second)
			}
			for _, pipe := range self.pipe {
				pipe.Stop()
			}
			return
		}
		for _, m := range self.matrices {
			if self.threads.JobCount() > int64(self.task.ThreadNum*2) {
				time.Sleep(time.Second)
			}
			if req := m.Pull(); req != nil {
				if ctx, err := self.down.Download(m.sp, req); err != nil {
					m.DoHistory(req, false)
				} else {
					// 过程处理，提炼数据
					ctx.Parse(req.GetRuleName())
					pipe := self.pipe[m.sp.GetName()]

					// 该条请求文件结果存入pipeline
					for _, f := range ctx.PullFiles() {
						if pipe.CollectFile(f) != nil {
							break
						}
					}
					// 该条请求文本结果存入pipeline
					for _, item := range ctx.PullItems() {
						if pipe.CollectData(item) != nil {
							break
						}
					}

					// 处理成功请求记录
					m.DoHistory(req, true)

					// 释放ctx准备复用
					spider.PutContext(ctx)

					sleeptime := 100
					time.Sleep(time.Duration(sleeptime) * time.Millisecond)
				}
			}
		}
	}
}

// 终止任务
func (self *Scheduler) Stop() {
	atomic.StoreInt32(&self.status, status.STOP)
}
