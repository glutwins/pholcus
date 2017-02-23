package history

import (
	"sync"

	"github.com/glutwins/pholcus/store"
)

type Success struct {
	tabName     string
	fileName    string
	new         map[string]bool // [Request.Unique()]true
	old         map[string]bool // [Request.Unique()]true
	inheritable bool
	sync.RWMutex
}

func (self *Success) Empty() {
	self.RWMutex.Lock()
	self.new = make(map[string]bool)
	self.old = make(map[string]bool)
	self.RWMutex.Unlock()
}

// 更新或加入成功记录，
// 对比是否已存在，不存在就记录，
// 返回值表示是否有插入操作。
func (self *Success) UpsertSuccess(reqUnique string) bool {
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()

	if self.old[reqUnique] {
		return false
	}
	if self.new[reqUnique] {
		return false
	}
	self.new[reqUnique] = true
	return true
}

func (self *Success) HasSuccess(reqUnique string) bool {
	self.RWMutex.Lock()
	has := self.old[reqUnique] || self.new[reqUnique]
	self.RWMutex.Unlock()
	return has
}

func (self *Success) flush(s store.Storage) (int, error) {
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()

	if len(self.new) == 0 {
		return 0, nil
	}

	docs := make(map[string]interface{})
	for k := range self.new {
		docs[k] = true
		self.old[k] = true
	}

	return s.InsertKVData(self.tabName, docs)
}
