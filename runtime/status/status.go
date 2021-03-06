package status

// 运行模式
const (
	UNSET int = iota - 1
	SERVER
	CLIENT
)

// 数据头部信息
const (
	// 任务请求Header
	REQTASK = iota + 1
	// 任务响应流头Header
	TASK
	// 打印Header
	LOG
)

const (
	PAUSE = iota
	STOP
	RUN
	STOPPED
)
