package controller

import (
	taskinteractor "algotrade_service/internal/engine/model/task"
)

type Controller struct {
	taskDispatcher taskinteractor.Dispatcher
	taskHandler    taskinteractor.Handler
	stop           chan struct{}
}

func NewController(taskDispatcher taskinteractor.Dispatcher, taskHandler taskinteractor.Handler) (*Controller, error) {
	return &Controller{
		taskDispatcher: taskDispatcher,
		taskHandler:    taskHandler,
	}, nil
}

func (c *Controller) Start() {
	c.stop = make(chan struct{})
	go c.runTaskDispatcher()
	go c.runTaskHandler()
}

func (c *Controller) Stop() {
	c.stop <- struct{}{}
}

func (c *Controller) runTaskDispatcher() {
	chTaskDispatcher := c.taskDispatcher.Subscribe()
	for {
		select {
		case task := <-chTaskDispatcher:
			c.taskHandler.Publish(task)
		case <-c.stop:
			return
		}
	}
}

func (c *Controller) runTaskHandler() {
	chTaskDispatcher := c.taskHandler.Subscribe()
	for {
		select {
		case resp := <-chTaskDispatcher:
			c.taskDispatcher.Publish(resp)
		case <-c.stop:
			return
		}
	}
}
