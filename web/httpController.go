package web

import (
	"net/http"
	"text/template"

	"github.com/glutwins/pholcus/app"
	"github.com/glutwins/pholcus/common/session"
	"github.com/glutwins/pholcus/config"
	"github.com/glutwins/pholcus/logs"
	"github.com/glutwins/pholcus/runtime/status"
)

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewManager("memory", `{"cookieName":"pholcusSession", "enableSetCookie,omitempty": true, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 157680000, "providerConfig": ""}`)
	// go globalSessions.GC()
}

// 处理web页面请求
func web(rw http.ResponseWriter, req *http.Request) {
	sess, _ := globalSessions.SessionStart(rw, req)
	defer sess.SessionRelease(rw)
	index, _ := viewsIndexHtmlBytes()
	t, err := template.New("index").Parse(string(index)) //解析模板文件
	// t, err := template.ParseFiles("web/views/index.html") //解析模板文件
	if err != nil {
		logs.Log.Error("%v", err)
	}
	//获取pholcus信息
	data := map[string]interface{}{
		"title":   config.NAME,
		"logo":    config.ICON_PNG,
		"version": config.VERSION,
		"author":  config.AUTHOR,
		"status": map[string]int{
			"stopped": status.STOPPED,
			"stop":    status.STOP,
			"run":     status.RUN,
			"pause":   status.PAUSE,
		},
		"port": app.LogicApp.GetAppConf("port").(int),
		"ip":   app.LogicApp.GetAppConf("master").(string),
	}
	t.Execute(rw, data) //执行模板的merger操作
}
