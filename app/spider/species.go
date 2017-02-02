package spider

// 全局蜘蛛种类实例

type SpiderSpecies map[string]*Spider

var Species = make(SpiderSpecies)

// 向蜘蛛种类清单添加新种类
func (self SpiderSpecies) Add(sp *Spider) {
	Species[sp.Name] = sp
}

func (self SpiderSpecies) GetByName(name string) *Spider {
	return self[name]
}
