package config

// 配置文件涉及的默认配置。
const (
	phantomjs               string = WORK_ROOT + "/phantomjs"    // phantomjs文件路径
	proxylib                string = WORK_ROOT + "/proxy.lib"    // 代理ip文件路径
	spiderdir               string = WORK_ROOT + "/spiders"      // 动态规则目录
	fileoutdir              string = WORK_ROOT + "/file_out"     // 文件（图片、HTML等）结果的输出目录
	textoutdir              string = WORK_ROOT + "/text_out"     // excel或csv输出方式下，文本结果的输出目录
	dbname                  string = TAG                         // 数据库名称
	mgoconnstring           string = "127.0.0.1:27017"           // mongodb连接字符串
	mgoconncap              int    = 1024                        // mongodb连接池容量
	mgoconngcsecond         int64  = 600                         // mongodb连接池GC时间，单位秒
	mysqlconnstring         string = "root:@tcp(127.0.0.1:3306)" // mysql连接字符串
	mysqlconncap            int    = 2048                        // mysql连接池容量
	mysqlmaxallowedpacketmb int    = 1                           //mysql通信缓冲区的最大长度，单位MB，默认1MB
	kafkabrokers            string = "127.0.0.1:9092"            //kafka broker字符串,逗号分割

	port        int    = 2015        // 主节点端口
	master      string = "127.0.0.1" // 服务器(主节点)地址，不含端口
	thread      int    = 20          // 全局最大并发量
	pause       int64  = 300         // 暂停时长参考/ms(随机: Pausetime/2 ~ Pausetime*2)
	outtype     string = "csv"       // 输出方式
	dockercap   int    = 10000       // 分段转储容器容量
	limit       int64  = 0           // 采集上限，0为不限，若在规则中设置初始值为LIMIT则为自定义限制，否则默认限制请求数
	proxyminute int64  = 0           // 代理IP更换的间隔分钟数
	success     bool   = true        // 继承历史成功记录
	failure     bool   = true        // 继承历史失败记录
)
