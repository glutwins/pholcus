// 数据收集
package pipeline

import (
	"github.com/glutwins/pholcus/app/pipeline/collector"
	"github.com/glutwins/pholcus/app/pipeline/collector/data"
	"github.com/glutwins/pholcus/app/spider"
	"github.com/glutwins/pholcus/common/schema"
)

// 数据收集/输出管道
type Pipeline interface {
	Start()                          //启动
	Stop()                           //停止
	CollectData(data.DataCell) error //收集数据单元
	CollectFile(data.FileCell) error //收集文件
}

func New(sp *spider.Spider, t *schema.Task) Pipeline {
	return collector.NewCollector(sp, t)
}
