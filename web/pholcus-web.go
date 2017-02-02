// [spider frame (golang)] Pholcus（幽灵蛛）是一款纯Go语言编写的高并发、分布式、重量级爬虫软件，支持单机、服务端、客户端三种运行模式，拥有Web、GUI、命令行三种操作界面；规则简单灵活、批量任务并发、输出方式丰富（mysql/mongodb/csv/excel等）、有大量Demo共享；同时她还支持横纵向两种抓取模式，支持模拟登录和任务暂停、取消等一系列高级功能；
//（官方QQ群：Go大数据 42731170，欢迎加入我们的讨论）。
// Web 界面版。
package web

import (
	"log"
	"mime"
	"net/http"
	"runtime"

	"github.com/glutwins/pholcus/app/spider"
	ws "github.com/glutwins/pholcus/common/websocket"
	"github.com/glutwins/pholcus/config"
)

var spiderMenu []map[string]string

// 执行入口
func Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Printf("%v\n\n", config.FULL_NAME)

	for name, sp := range spider.Species {
		spiderMenu = append(spiderMenu, map[string]string{"name": name, "description": sp.Description})
	}

	mime.AddExtensionType(".css", "text/css")

	// 设置websocket请求路由
	http.Handle("/ws", ws.Handler(wsHandle))
	// 设置websocket报告打印专用路由
	http.Handle("/ws/log", ws.Handler(wsLogHandle))
	//设置http访问的路由
	http.HandleFunc("/", web)
	//static file server
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(assetFS())))

	log.Printf("[pholcus] Server running on %v\n", config.DefaultConfig.WebAddr)
	// 监听端口
	if err := http.ListenAndServe(config.DefaultConfig.WebAddr, nil); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
