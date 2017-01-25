// [spider frame (golang)] Pholcus（幽灵蛛）是一款纯Go语言编写的高并发、分布式、重量级爬虫软件，支持单机、服务端、客户端三种运行模式，拥有Web、GUI、命令行三种操作界面；规则简单灵活、批量任务并发、输出方式丰富（mysql/mongodb/csv/excel等）、有大量Demo共享；同时她还支持横纵向两种抓取模式，支持模拟登录和任务暂停、取消等一系列高级功能；
//（官方QQ群：Go大数据 42731170，欢迎加入我们的讨论）。
// Web 界面版。
package web

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/glutwins/pholcus/app"
	"github.com/glutwins/pholcus/config"
	"github.com/glutwins/pholcus/logs"
	"github.com/glutwins/pholcus/runtime/cache"
)

var (
	spiderMenu []map[string]string
	keyinsflag *string
	limitflag  *int64
	threadflag *int
	pauseflag  *int64
	proxyflag  *int64
	dockerflag *int
)

// 执行入口
func Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Printf("%v\n\n", config.FULL_NAME)

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

	flag.String("z", "", "README:   参数设置参考 [xxx] 提示，参数中包含多个值时以 \",\" 间隔。\r\n")
	flag.Parse()

	cache.Task.Keyins = *keyinsflag
	cache.Task.Limit = *limitflag
	cache.Task.ThreadNum = *threadflag
	cache.Task.Pausetime = *pauseflag
	cache.Task.ProxyMinute = *proxyflag
	cache.Task.DockerCap = *dockerflag

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
}
