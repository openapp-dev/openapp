package utils

import (
	pkgtypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
)

type WorkItemHandler func(pkgtypes.NamespacedName) error

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

func (w *WorkQueue) Add(item pkgtypes.NamespacedName) {
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

	namespaceNameKey := obj.(pkgtypes.NamespacedName)
	if err := w.handler(namespaceNameKey); err != nil {
		w.queue.AddRateLimited(obj)
		return true
	}

	return true
}
