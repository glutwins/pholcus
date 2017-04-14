package slave

import (
	"context"
	"time"

	"github.com/glutwins/pholcus/app/spider"
	"github.com/glutwins/pholcus/common/schema"
)

type TaskProcess struct {
	id     int
	cancel context.CancelFunc // 停止整个任务
}

func (tp *TaskProcess) Run(ctx context.Context, task *schema.Task) (*schema.TaskResult, error) {
	ch := make(chan *schema.TaskResult)
	go tp.Schedule(ctx, task, ch)

	select {
	case <-ctx.Done():
		return <-ch, ctx.Err()
	case result := <-ch:
		return result, nil
	}
}

func (tp *TaskProcess) Schedule(ctx context.Context, task *schema.Task, ch chan *schema.TaskResult) {
	result := &schema.TaskResult{}
	result.StartTime = time.Now()

	var matrixs []*Matrix
	for name, keyins := range task.Spiders {
		sp := spider.Species.GetByName(name)
		sp.RuleTree.Root(spider.GetContext(sp, nil))
		for _, keyin := range keyins {
			m := NewMatrix(name, keyin, task, sp)
			matrixs = append(matrixs, m)

			go m.Run(ctx)
		}
	}

	for {
		select {
		case <-ctx.Done():
			break
		case <-time.After(time.Millisecond * 10):
			finish := true
			for i := 0; i < len(matrixs); i++ {

			}

			if finish {
				break
			}
		}
	}

	result.TakeTime = time.Since(result.StartTime)
	ch <- result
}
