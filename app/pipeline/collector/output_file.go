package collector

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/glutwins/pholcus/app/pipeline/collector/data"
	bytesSize "github.com/glutwins/pholcus/common/bytes"
	"github.com/glutwins/pholcus/common/util"
	"github.com/glutwins/pholcus/config"
	"github.com/glutwins/pholcus/logs"
)

// 文件输出
func (self *Collector) outputFile(file data.FileCell) {
	// 复用FileCell
	defer func() {
		data.PutFileCell(file)
		self.wait.Done()
	}()

	p, n := filepath.Split(filepath.Clean(file["Name"].(string)))
	dir := filepath.Join(config.DefaultConfig.FileOutDir, self.name, p)

	// 文件名
	fileName := filepath.Join(dir, util.FileNameReplace(n))

	// 创建/打开目录
	d, err := os.Stat(dir)
	if err != nil || !d.IsDir() {
		if err := os.MkdirAll(dir, 0777); err != nil {
			logs.Log.Error(" *     Fail  [文件下载： | KEYIN： | 批次：%v]   %v [ERROR]  %v\n",
				atomic.LoadUint64(&self.fileBatch), fileName, err,
			)
			return
		}
	}

	// 文件不存在就以0777的权限创建文件，如果存在就在写入之前清空内容
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		logs.Log.Error(
			" *     Fail  [文件下载： | KEYIN： | 批次：%v]   %v [ERROR]  %v\n",
			atomic.LoadUint64(&self.fileBatch), fileName, err,
		)
		return
	}

	size, err := io.Copy(f, bytes.NewReader(file["Bytes"].([]byte)))
	f.Close()
	if err != nil {
		logs.Log.Error(
			" *     Fail  [文件下载： | KEYIN： | 批次：%v]   %v (%s) [ERROR]  %v\n",
			atomic.LoadUint64(&self.fileBatch), fileName, bytesSize.Format(uint64(size)), err,
		)
		return
	}

	// 输出统计
	self.addFileSum(1)

	// 打印报告
	logs.Log.Informational(" * ")
	logs.Log.Informational(
		" *     [文件下载：| KEYIN：| 批次：%v]   %v (%s)\n",
		atomic.LoadUint64(&self.fileBatch), fileName, bytesSize.Format(uint64(size)),
	)
	logs.Log.Informational(" * ")
}
