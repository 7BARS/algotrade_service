package task

import "algotrade_service/internal/model/task"

type Handler interface {
	Subscribe() <- chan task.Response
	Publish(task.Task)
}