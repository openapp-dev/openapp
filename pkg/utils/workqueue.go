package utils

import (
	"k8s.io/client-go/util/workqueue"
)

type WorkItemHandler func(obj interface{}) error

type WorkQueue struct {
	queue   workqueue.RateLimitingInterface
	handler WorkItemHandler
}

func NewWorkQueue(handler WorkItemHandler) *WorkQueue {
	return &WorkQueue{
		queue:   workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		handler: handler,
	}
}

func (w *WorkQueue) Add(item interface{}) {
	w.queue.AddRateLimited(item)
}

func (w *WorkQueue) Run() {
	for w.process() {
	}
}

func (w *WorkQueue) process() bool {
	obj, quit := w.queue.Get()
	if quit {
		return false
	}
	defer w.queue.Done(obj)

	if err := w.handler(obj); err != nil {
		w.queue.AddRateLimited(obj)
		return true
	}

	return true
}
