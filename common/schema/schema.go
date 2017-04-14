package schema

import (
	"time"

	"github.com/glutwins/pholcus/config"
)

type Task struct {
	Id             int
	Spiders        map[string][]string // 蜘蛛规则name字段与keyin字段，规定格式map[string]string{"name":"baidu","keyin":"henry"}
	ThreadNum      int                 // 全局最大并发量
	Pausetime      int64               // 暂停时长参考/ms(随机: Pausetime/2 ~ Pausetime*2)
	DockerCap      int                 // 分段转储容器容量
	DockerQueueCap int                 // 分段输出池容量，不小于2
	Inherit        bool                // 继承历史记录
	UseProxy       bool                // 是否使用代理
	Limit          int64               // 采集上限，0为不限，若在规则中设置初始值为LIMIT则为自定义限制，否则默认限制请求数
	// 选填项
	Keyins string // 自定义输入，后期切分为多个任务的Keyin自定义配置

	Db *config.PholcusDbConfig // 定义具体的采集规则树

}

type TaskResult struct {
	StartTime time.Time
	TakeTime  time.Duration // 花费时长
	DataNum   int           // 处理数据量
	FileNum   int           // 处理文件量
	FailNum   int           // 失败数量
	SuccNum   int           // 成功数量
}
