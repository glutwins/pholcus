package history

import (
	"sync"

	"github.com/glutwins/pholcus/app/downloader/request"
	"github.com/glutwins/pholcus/store"
)

type Failure struct {
	tabName     string
	list        map[string]*request.Request //key:url
	inheritable bool
	sync.RWMutex
}

func (self *Failure) Empty() {
	self.RWMutex.Lock()
	self.list = make(map[string]*request.Request)
	self.RWMutex.Unlock()
}

func (self *Failure) PullFailure() map[string]*request.Request {
	list := self.list
	self.list = make(map[string]*request.Request)
	return list
}

// 更新或加入失败记录，
// 对比是否已存在，不存在就记录，
// 返回值表示是否有插入操作。
func (self *Failure) UpsertFailure(req *request.Request) bool {
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()
	if self.list[req.Unique()] != nil {
		return false
	}
	self.list[req.Unique()] = req
	return true
}

// 先清空历史失败记录再更新
func (self *Failure) flush(w store.Storage) (int, error) {
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()

	docs := make(map[string]interface{}, len(self.list))
	for key, req := range self.list {
		docs[key] = req.Serialize()
	}

	w.ClearKVData(self.tabName)
	return w.InsertKVData(self.tabName, docs)
}
