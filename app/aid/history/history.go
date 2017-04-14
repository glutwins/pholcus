package history

import (
	"sync"
	"time"

	"github.com/glutwins/pholcus/app/downloader/request"
	"github.com/glutwins/pholcus/config"
	"github.com/glutwins/pholcus/logs"
	"github.com/glutwins/pholcus/store"
)

type historyValue struct {
	Succ  bool
	Url   string
	Exval string
}

type History struct {
	succmap map[string]*historyValue
	failmap map[string]*historyValue
	tabName string
	s       store.Storage
	inherit bool
	sync.RWMutex
}

func New(name string, subName string, db *config.PholcusDbConfig) *History {
	self := &History{
		succmap: make(map[string]*historyValue),
		failmap: make(map[string]*historyValue),
		s:       store.NewStorage(db),
	}

	self.tabName = config.HISTORY_TAG + "__y__" + name
	if subName != "" {
		self.tabName += "__" + subName
	}

	if self.inherit {
		docs, _ := self.s.FetchKVData(self.tabName)
		for k, v := range docs {
			r := v.(*historyValue)
			if r.Succ {
				self.succmap[k] = r
			} else {
				self.failmap[k] = r
			}
		}

		logs.Log.Informational(" *     [读取记录]: 成功 %d 条, 失败 %d 条\n", len(self.succmap), len(self.failmap))
	}

	return self
}

func (self *History) Upsert(req *request.Request, err error) {
	self.Lock()
	defer self.Unlock()

	uid := req.Unique()
	if _, ok := self.succmap[uid]; !ok {
		if _, ok := self.failmap[uid]; ok {
			if err == nil {
				delete(self.failmap, uid)
				self.succmap[uid] = &historyValue{true, req.Url, time.Now().String()}
			}
		} else {
			v := &historyValue{Url: req.Url}
			if err != nil {
				v.Exval = err.Error()
				self.failmap[uid] = v
			} else {
				v.Succ = true
				v.Exval = time.Now().String()
				self.succmap[uid] = v
			}
		}
	}
}

func (self *History) HasSucc(uid string) bool {
	return true
}

// I/O输出成功记录，但不清缓存
func (self *History) Flush() {
	self.Lock()
	defer self.Unlock()

	result := make(map[string]interface{})
	for k, v := range self.succmap {
		result[k] = v
	}
	for k, v := range self.failmap {
		result[k] = v
	}

	if len(result) > 0 {
		self.s.InsertKVData(self.tabName, result)
	}
}
