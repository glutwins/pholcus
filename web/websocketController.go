package web

import (
	"sync"

	"github.com/glutwins/pholcus/app/spider"
	"github.com/glutwins/pholcus/common/schema"
	"github.com/glutwins/pholcus/common/util"
	ws "github.com/glutwins/pholcus/common/websocket"
	"github.com/glutwins/pholcus/config"
	"github.com/glutwins/pholcus/logs"
	"github.com/glutwins/pholcus/runtime/status"
)

type SocketController struct {
	connPool     map[string]*ws.Conn
	wchanPool    map[string]*Wchan
	connRWMutex  sync.RWMutex
	wchanRWMutex sync.RWMutex
}

func (self *SocketController) GetConn(sessID string) *ws.Conn {
	self.connRWMutex.RLock()
	defer self.connRWMutex.RUnlock()
	return self.connPool[sessID]
}

func (self *SocketController) GetWchan(sessID string) *Wchan {
	self.wchanRWMutex.RLock()
	defer self.wchanRWMutex.RUnlock()
	return self.wchanPool[sessID]
}

func (self *SocketController) Add(sessID string, conn *ws.Conn) {
	self.connRWMutex.Lock()
	self.wchanRWMutex.Lock()
	defer self.connRWMutex.Unlock()
	defer self.wchanRWMutex.Unlock()

	self.connPool[sessID] = conn
	self.wchanPool[sessID] = newWchan()
}

func (self *SocketController) Remove(sessID string, conn *ws.Conn) {
	self.connRWMutex.Lock()
	self.wchanRWMutex.Lock()
	defer self.connRWMutex.Unlock()
	defer self.wchanRWMutex.Unlock()

	if self.connPool[sessID] == nil {
		return
	}
	wc := self.wchanPool[sessID]
	close(wc.wchan)
	conn.Close()
	delete(self.connPool, sessID)
	delete(self.wchanPool, sessID)
}

func (self *SocketController) Write(sessID string, void map[string]interface{}, to ...int) {
	self.wchanRWMutex.RLock()
	defer self.wchanRWMutex.RUnlock()

	// to为1时，只向当前连接发送；to为-1时，向除当前连接外的其他所有连接发送；to为0时或为空时，向所有连接发送
	var t int = 0
	if len(to) > 0 {
		t = to[0]
	}

	switch t {
	case 1:
		wc := self.wchanPool[sessID]
		if wc == nil {
			return
		}
		void["initiative"] = true
		wc.wchan <- void

	case 0, -1:
		l := len(self.wchanPool)
		for _sessID, wc := range self.wchanPool {
			if t == -1 && _sessID == sessID {
				continue
			}
			_void := make(map[string]interface{}, l)
			for k, v := range void {
				_void[k] = v
			}
			if _sessID == sessID {
				_void["initiative"] = true
			} else {
				_void["initiative"] = false
			}
			wc.wchan <- _void
		}
	}
}

type Wchan struct {
	wchan chan interface{}
}

func newWchan() *Wchan {
	return &Wchan{
		wchan: make(chan interface{}, 1024),
	}
}

var (
	wsApi = map[string]func(string, map[string]interface{}){}
	Sc    = &SocketController{
		connPool:  make(map[string]*ws.Conn),
		wchanPool: make(map[string]*Wchan),
	}
)

func wsHandle(conn *ws.Conn) {
	defer func() {
		if p := recover(); p != nil {
			logs.Log.Error("%v", p)
		}
	}()
	sess, _ := globalSessions.SessionStart(nil, conn.Request())
	sessID := sess.SessionID()
	if Sc.GetConn(sessID) == nil {
		Sc.Add(sessID, conn)
	}

	defer Sc.Remove(sessID, conn)

	go func() {
		var err error
		for info := range Sc.GetWchan(sessID).wchan {
			if _, err = ws.JSON.Send(conn, info); err != nil {
				return
			}
		}
	}()

	for {
		var req map[string]interface{}

		if err := ws.JSON.Receive(conn, &req); err != nil {
			return
		}

		wsApi[util.Atoa(req["operate"])](sessID, req)
	}
}

func init() {

	// 初始化运行
	wsApi["refresh"] = func(sessID string, req map[string]interface{}) {
		// 写入发送通道
		Sc.Write(sessID, tplData(status.SERVER), 1)
	}

	// 初始化运行
	wsApi["init"] = func(sessID string, req map[string]interface{}) {
		logs.Log.SetOutput(Lsc)
		// 写入发送通道
		Sc.Write(sessID, tplData(status.SERVER))
	}

	wsApi["run"] = func(sessID string, req map[string]interface{}) {
		t := &schema.Task{}
		t.ThreadNum = util.Atoi(req["ThreadNum"])
		t.Pausetime = int64(util.Atoi(req["Pausetime"]))
		t.ProxyMinute = int64(util.Atoi(req["ProxyMinute"]))
		t.DockerCap = util.Atoi(req["DockerCap"])
		t.Limit = int64(util.Atoi(req["Limit"]))
		t.Keyins = util.Atoa(req["Keyins"])
		t.Inherit = req["Inherit"] == "true"

		spNames, ok := req["spiders"].([]interface{})
		if !ok {
			return
		}
		spiders := []*spider.Spider{}
		for _, sp := range app.LogicApp.GetSpiderLib() {
			for _, spName := range spNames {
				if util.Atoa(spName) == sp.GetName() {
					spiders = append(spiders, sp.Copy())
				}
			}
		}
		app.LogicApp.SpiderPrepare(spiders)
	}

	// 终止当前任务，现仅支持单机模式
	wsApi["stop"] = func(sessID string, req map[string]interface{}) {
		Sc.Write(sessID, map[string]interface{}{"operate": "stop"})

	}

	// 任务暂停与恢复，目前仅支持单机模式
	wsApi["pauseRecover"] = func(sessID string, req map[string]interface{}) {
	}
}

func tplData(mode int) map[string]interface{} {
	var info = map[string]interface{}{"operate": "init", "mode": mode}

	// 运行模式标题
	info["title"] = config.FULL_NAME + "                                                          【 运行模式 ->  服务端 】"

	// 蜘蛛家族清单
	info["spiders"] = map[string]interface{}{
		"menu": spiderMenu,
		"curr": func() interface{} {
			l := app.LogicApp.GetSpiderQueue().Len()
			if l == 0 {
				return 0
			}
			var curr = make(map[string]bool, l)
			for _, sp := range app.LogicApp.GetSpiderQueue().GetAll() {
				curr[sp.GetName()] = true
			}

			return curr
		}(),
	}

	// 输出方式清单
	info["OutType"] = map[string]interface{}{
		"menu": config.DefaultConfig.Db,
	}

	// 并发协程上限
	info["ThreadNum"] = map[string]int{
		"max":  999999,
		"min":  1,
		"curr": 1,
	}

	// 暂停区间/ms(随机: Pausetime/2 ~ Pausetime*2)
	info["Pausetime"] = map[string][]int64{
		"menu": {0, 100, 300, 500, 1000, 3000, 5000, 10000, 15000, 20000, 30000, 60000},
		"curr": []int64{100},
	}

	// 代理IP更换的间隔分钟数
	info["ProxyMinute"] = map[string][]int64{
		"menu": {0, 1, 3, 5, 10, 15, 20, 30, 45, 60, 120, 180},
		"curr": []int64{1},
	}

	// 分批输出的容量
	info["DockerCap"] = map[string]int{
		"min":  1,
		"max":  5000000,
		"curr": 1,
	}

	// 采集上限
	info["Limit"] = 0

	// 自定义配置
	info["Keyins"] = ""

	return info
}
