package service

import (
	"fmt"

	"algotrade_service/internal/engine/model/data"
	"algotrade_service/internal/model/task"
)

type TaskManager struct {
	task              []task.Task
	chTaskResponse    <-chan task.Response
	dataHistoryChange data.HistoryChange
	chStop            chan struct{}
	// expressions
}

func NewService(dataHistoryChange data.HistoryChange) *TaskManager {
	return &TaskManager{
		dataHistoryChange: dataHistoryChange,
	}
}

func (s *TaskManager) Start() {
	s.task = make([]task.Task, 0)
	s.chStop = make(chan struct{})
	s.chTaskResponse = make(<-chan task.Response)
	go s.run()
}

func (s *TaskManager) Subscribe() <-chan task.Response {
	return s.chTaskResponse
}

func (s *TaskManager) Publish(task task.Task) {
	// s.task = append(s.task, task)
	s.processMsg(task)
}

func (s *TaskManager) run() {
	chTicker := s.dataHistoryChange.Subscribe()
	for {
		select {
		case ticker := <-chTicker:
			fmt.Printf("ticker: %v\n", ticker)
		case <-s.chStop:
			return
		}
	}
}

// processMsg()
func (s *TaskManager) processMsg(task task.Task) {
	
}