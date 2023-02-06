package task

import (
	"algotrade_service/internal/model/task"
)

type Dispatcher interface {
	Subscribe() <- chan task.Task
	Publish(task.Response)
}
