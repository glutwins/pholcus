// [spider frame (golang)] Pholcus（幽灵蛛）是一款纯Go语言编写的高并发、分布式、重量级爬虫软件，支持单机、服务端、客户端三种运行模式，拥有Web、GUI、命令行三种操作界面；规则简单灵活、批量任务并发、输出方式丰富（mysql/mongodb/csv/excel等）、有大量Demo共享；同时她还支持横纵向两种抓取模式，支持模拟登录和任务暂停、取消等一系列高级功能；
//（官方QQ群：Go大数据 42731170，欢迎加入我们的讨论）。
// Web 界面版。
package web

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/glutwins/pholcus/app"
	"github.com/glutwins/pholcus/common/gc"
	"github.com/glutwins/pholcus/config"
	"github.com/glutwins/pholcus/logs"
	"github.com/glutwins/pholcus/runtime/cache"
)

var (
	spiderMenu         []map[string]string
	portflag           *int
	masterflag         *string
	keyinsflag         *string
	limitflag          *int64
	threadflag         *int
	pauseflag          *int64
	proxyflag          *int64
	dockerflag         *int
	successInheritflag *bool
	failureInheritflag *bool
)

// 执行入口
func Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	gc.ManualGC()

	fmt.Printf("%v\n\n", config.FULL_NAME)
	flag.String("a *********************************************** common *********************************************** -a", "", "")
	// 操作界面
	//端口号，非单机模式填写
	portflag = flag.Int(
		"a_port",
		cache.Task.Port,
		"   <端口号: 只填写数字即可，不含冒号，单机模式不填>")

	//主节点ip，客户端模式填写
	masterflag = flag.String(
		"a_master",
		cache.Task.Master,
		"   <服务端IP: 不含端口，客户端模式下使用>")

	// 自定义配置
	keyinsflag = flag.String(
		"a_keyins",
		cache.Task.Keyins,
		"   <自定义配置: 多任务请分别多包一层“<>”>")

	// 采集上限
	limitflag = flag.Int64(
		"a_limit",
		cache.Task.Limit,
		"   <采集上限（默认限制URL数）> [>=0]")

	// 并发协程数
	threadflag = flag.Int(
		"a_thread",
		cache.Task.ThreadNum,
		"   <并发协程> [1~99999]")

	// 平均暂停时间
	pauseflag = flag.Int64(
		"a_pause",
		cache.Task.Pausetime,
		"   <平均暂停时间/ms> [>=100]")

	// 代理IP更换频率
	proxyflag = flag.Int64(
		"a_proxyminute",
		cache.Task.ProxyMinute,
		"   <代理IP更换频率: /m，为0时不使用代理> [>=0]")

	// 分批输出
	dockerflag = flag.Int(
		"a_dockercap",
		cache.Task.DockerCap,
		"   <分批输出> [1~5000000]")

	// 继承历史成功记录
	successInheritflag = flag.Bool(
		"a_success",
		cache.Task.SuccessInherit,
		"   <继承并保存成功记录> [true] [false]")

	// 继承历史失败记录
	failureInheritflag = flag.Bool(
		"a_failure",
		cache.Task.FailureInherit,
		"   <继承并保存失败记录> [true] [false]")

	flag.String("z", "", "README:   参数设置参考 [xxx] 提示，参数中包含多个值时以 \",\" 间隔。\r\n")
	flag.Parse()

	cache.Task.Port = *portflag
	cache.Task.Master = *masterflag
	cache.Task.Keyins = *keyinsflag
	cache.Task.Limit = *limitflag
	cache.Task.ThreadNum = *threadflag
	cache.Task.Pausetime = *pauseflag
	cache.Task.ProxyMinute = *proxyflag
	cache.Task.DockerCap = *dockerflag
	cache.Task.SuccessInherit = *successInheritflag
	cache.Task.FailureInherit = *failureInheritflag

	ctrl := make(chan os.Signal, 1)
	signal.Notify(ctrl, os.Interrupt, os.Kill)

	go func() {
		spiderMenu = func() (spmenu []map[string]string) {
			// 获取蜘蛛家族
			for _, sp := range app.LogicApp.GetSpiderLib() {
				spmenu = append(spmenu, map[string]string{"name": sp.GetName(), "description": sp.GetDescription()})
			}
			return spmenu
		}()

		// 预绑定路由
		Router()

		log.Printf("[pholcus] Server running on %v\n", config.DefaultConfig.WebAddr)
		// 监听端口
		if err := http.ListenAndServe(config.DefaultConfig.WebAddr, nil); err != nil {
			logs.Log.Error("ListenAndServe: %v", err)
		}
	}()

	<-ctrl
}
