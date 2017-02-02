package history

import (
	"github.com/glutwins/pholcus/app/downloader/request"
	"github.com/glutwins/pholcus/common/util"
	"github.com/glutwins/pholcus/config"
	"github.com/glutwins/pholcus/logs"
	"github.com/glutwins/pholcus/store"
)

type (
	Historier interface {
		ReadSuccess()              // 读取成功记录
		UpsertSuccess(string) bool // 更新或加入成功记录
		HasSuccess(string) bool    // 检查是否存在某条成功记录
		DeleteSuccess(string)      // 删除成功记录
		FlushSuccess()             // I/O输出成功记录，但不清缓存

		ReadFailure()                             // 取出失败记录
		PullFailure() map[string]*request.Request // 拉取失败记录并清空
		UpsertFailure(*request.Request) bool      // 更新或加入失败记录
		DeleteFailure(*request.Request)           // 删除失败记录
		FlushFailure()                            // I/O输出失败记录，但不清缓存

		Empty() // 清空缓存，但不输出
	}
	History struct {
		*Success
		*Failure
		s       store.Storage
		inherit bool
	}
)

func New(name string, subName string, db *config.PholcusDbConfig) Historier {
	successTabName := config.HISTORY_TAG + "__y__" + name
	failureTabName := config.HISTORY_TAG + "__n__" + name
	if subName != "" {
		successTabName += "__" + subName
		failureTabName += "__" + subName
	}
	return &History{
		Success: &Success{
			tabName: util.FileNameReplace(successTabName),
			new:     make(map[string]bool),
			old:     make(map[string]bool),
		},
		Failure: &Failure{
			tabName: util.FileNameReplace(failureTabName),
			list:    make(map[string]*request.Request),
		},
		s: store.NewStorage(db),
	}
}

// 读取成功记录
func (self *History) ReadSuccess() {
	if !self.inherit {
		// 不继承历史记录时
		self.Success.old = make(map[string]bool)
		self.Success.new = make(map[string]bool)
		self.Success.inheritable = false
		return

	} else if self.Success.inheritable {
		// 本次与上次均继承历史记录时
		return

	} else {
		// 上次没有继承历史记录，但本次继承时
		self.Success.old = make(map[string]bool)
		self.Success.new = make(map[string]bool)
		self.Success.inheritable = true
	}

	docs, _ := self.s.FetchKVData(self.Success.tabName)
	for k, _ := range docs {
		self.Success.old[k] = true
	}

	logs.Log.Informational(" *     [读取成功记录]: %v 条\n", len(self.Success.old))
}

// 取出失败记录
func (self *History) ReadFailure() {
	if !self.inherit {
		// 不继承历史记录时
		self.Failure.list = make(map[string]*request.Request)
		self.Failure.inheritable = false
		return

	} else if self.Failure.inheritable {
		// 本次与上次均继承历史记录时
		return

	} else {
		// 上次没有继承历史记录，但本次继承时
		self.Failure.list = make(map[string]*request.Request)
		self.Failure.inheritable = true
	}

	docs, _ := self.s.FetchKVData(self.Failure.tabName)
	for k, v := range docs {
		self.Failure.list[k], _ = request.UnSerialize(v.(string))
	}

	logs.Log.Informational(" *     [取出失败记录]: %v 条\n", len(docs))
}

// 清空缓存，但不输出
func (self *History) Empty() {
	self.Success.Empty()
	self.Failure.Empty()
}

// I/O输出成功记录，但不清缓存
func (self *History) FlushSuccess() {
	if !self.Success.inheritable {
		return
	}

	sucLen, err := self.Success.flush(self.s)
	if sucLen <= 0 {
		return
	}
	if err != nil {
		logs.Log.Error("%v", err)
	} else {
		logs.Log.Informational(" *     [添加成功记录]: %v 条\n", sucLen)
	}
}

// I/O输出失败记录，但不清缓存
func (self *History) FlushFailure() {
	if !self.Failure.inheritable {
		return
	}
	failLen, err := self.Failure.flush(self.s)
	if failLen <= 0 {
		return
	}
	if err != nil {
		logs.Log.Error("%v", err)
	} else {
		logs.Log.Informational(" *     [添加失败记录]: %v 条\n", failLen)
	}
}
