package spider

import (
	"github.com/glutwins/pholcus/common/util"
)

type (
	// 蜘蛛规则
	Spider struct {
		// 以下字段由用户定义
		Name         string                                       // 用户界面显示的名称（应保证唯一性）
		Description  string                                       // 用户界面显示的描述
		Pausetime    int64                                        // 随机暂停区间(50%~200%)，若规则中直接定义，则不被界面传参覆盖
		Limit        int64                                        // 默认限制请求数，0为不限；若规则中定义为LIMIT，则采用规则的自定义限制方案
		Keyin        string                                       // 自定义输入的配置信息，使用前须在规则中设置初始值为KEYIN
		EnableCookie bool                                         // 所有请求是否使用cookie记录
		Namespace    func(dataCell map[string]interface{}) string // 次级命名，用于输出文件、路径的命名，可依赖具体数据内容
		RuleTree     *RuleTree
	}
	//采集规则树
	RuleTree struct {
		Root  func(*Context)   // 根节点(执行入口)
		Trunk map[string]*Rule // 节点散列表(执行采集过程)
	}
	// 采集规则节点
	Rule struct {
		ParseFunc func(*Context)                                     // 内容解析函数
		AidFunc   func(*Context, map[string]interface{}) interface{} // 通用辅助函数
	}
)

// 添加自身到蜘蛛菜单
func (self Spider) Register() {
	Species.Add(&self)
}

// 获取蜘蛛名称
func (self *Spider) GetName() string {
	return self.Name
}

// 获取蜘蛛二级标识名
func (self *Spider) GetSubName() string {
	if len([]rune(self.Keyin)) > 8 {
		return util.MakeHash(self.Keyin)
	}

	return self.Keyin
}
