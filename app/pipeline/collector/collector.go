// 结果收集与输出
package collector

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/glutwins/pholcus/app/pipeline/collector/data"
	"github.com/glutwins/pholcus/app/spider"
	"github.com/glutwins/pholcus/common/schema"
	"github.com/glutwins/pholcus/common/util"
	"github.com/glutwins/pholcus/logs"
	"github.com/glutwins/pholcus/store"
)

// 结果收集与输出
type Collector struct {
	DataChan    chan data.DataCell //文本数据收集通道
	FileChan    chan data.FileCell //文件收集通道
	dataDocker  []data.DataCell    //分批输出结果缓存
	dataBatch   uint64             //当前文本输出批次
	fileBatch   uint64             //当前文件输出批次
	wait        sync.WaitGroup
	sum         [4]uint64 //收集的数据总数[上次输出后文本总数，本次输出后文本总数，上次输出后文件总数，本次输出后文件总数]，非并发安全
	dataSumLock sync.RWMutex
	fileSumLock sync.RWMutex

	name  string
	subns func(map[string]interface{}) string

	store store.Storage
}

func NewCollector(sp *spider.Spider, t *schema.Task) *Collector {
	var self = &Collector{}
	self.name = sp.GetSubName()
	if self.name == "" {
		self.name = sp.Name
	} else {
		self.name = sp.Name + "__" + self.name
	}
	self.name = util.FileNameReplace(self.name)

	if sp.Namespace != nil {
		self.subns = sp.Namespace
	} else {
		self.subns = func(cell map[string]interface{}) string {
			return cell["RuleName"].(string)
		}
	}

	self.DataChan = make(chan data.DataCell, t.DockerCap)
	self.FileChan = make(chan data.FileCell, t.DockerCap)
	self.dataDocker = make([]data.DataCell, 0, t.DockerCap)
	self.sum = [4]uint64{}
	self.dataBatch = 0
	self.fileBatch = 0
	self.store = store.NewStorage(t.Db)
	return self
}

func (self *Collector) CollectData(dataCell data.DataCell) error {
	var err error
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("输出协程已终止")
		}
	}()
	self.DataChan <- dataCell
	return err
}

func (self *Collector) CollectFile(fileCell data.FileCell) error {
	var err error
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("输出协程已终止")
		}
	}()
	self.FileChan <- fileCell
	return err
}

// 停止
func (self *Collector) Stop() {
	go func() {
		defer func() {
			recover()
		}()
		close(self.DataChan)
	}()
	go func() {
		defer func() {
			recover()
		}()
		close(self.FileChan)
	}()
}

// 启动数据收集/输出管道
func (self *Collector) Start() {
	// 启动输出协程
	go func() {
		dataStop := make(chan bool)
		fileStop := make(chan bool)

		go func() {
			defer func() {
				recover()
			}()
			for data := range self.DataChan {
				// 缓存分批数据
				self.dataDocker = append(self.dataDocker, data)

				// 未达到设定的分批量时继续收集数据
				if len(self.dataDocker) < cap(self.dataDocker) {
					continue
				}

				// 执行输出
				self.dataBatch++
				self.outputData()
			}
			// 将剩余收集到但未输出的数据输出
			self.dataBatch++
			self.outputData()
			close(dataStop)
		}()

		go func() {
			defer func() {
				recover()
			}()
			// 只有当收到退出通知并且通道内无数据时，才退出循环
			for file := range self.FileChan {
				atomic.AddUint64(&self.fileBatch, 1)
				self.wait.Add(1)
				go self.outputFile(file)
			}
			close(fileStop)
		}()

		<-dataStop
		<-fileStop

		// 等待所有输出完成
		self.wait.Wait()
	}()
}

// 获取文本数据总量
func (self *Collector) dataSum() uint64 {
	self.dataSumLock.RLock()
	defer self.dataSumLock.RUnlock()
	return self.sum[1]
}

// 更新文本数据总量
func (self *Collector) addDataSum(add uint64) {
	self.dataSumLock.Lock()
	defer self.dataSumLock.Unlock()
	self.sum[0] = self.sum[1]
	self.sum[1] += add
}

// 获取文件数据总量
func (self *Collector) fileSum() uint64 {
	self.fileSumLock.RLock()
	defer self.fileSumLock.RUnlock()
	return self.sum[3]
}

// 更新文件数据总量
func (self *Collector) addFileSum(add uint64) {
	self.fileSumLock.Lock()
	defer self.fileSumLock.Unlock()
	self.sum[2] = self.sum[3]
	self.sum[3] += add
}

// 文本数据输出
func (self *Collector) outputData() {
	defer func() {
		for _, cell := range self.dataDocker {
			data.PutDataCell(cell)
		}
		self.dataDocker = self.dataDocker[:0]
	}()

	// 输出
	dataLen := uint64(len(self.dataDocker))
	if dataLen == 0 {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			logs.Log.Informational(" *     Panic  %v\n", p)
		}
	}()

	// 输出统计
	self.addDataSum(dataLen)
	now := time.Now()

	// 执行输出
	for _, datacell := range self.dataDocker {
		ns := self.name
		if sub := self.subns(datacell); sub != "" {
			ns = ns + "__" + util.FileNameReplace(sub)
		}
		tblname := now.Format("2006-01-02 150405") + "/" + ns

		row := map[string]interface{}{}
		row["当前链接"] = datacell["Url"]
		row["上级链接"] = datacell["ParentUrl"]
		row["下载时间"] = datacell["DownloadTime"]
		vd := datacell["Data"].(map[string]interface{})
		for k, v := range vd {
			row[k] = v
		}

		if err := self.store.InsertStringMap(tblname, row); err != nil {
			logs.Log.Error(" *     Fail  [数据输出： | KEYIN： | 批次：%v]   数据 %v 条！ [ERROR]  %v\n",
				self.dataBatch, dataLen, err)
		}
	}

	logs.Log.Informational(" *     [数据输出：| KEYIN： | 批次：%v]   数据 %v 条！\n",
		self.dataBatch, dataLen)
	//self.Spider.TryFlushSuccess()
}
