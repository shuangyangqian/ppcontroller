package task

import (
	"k8s.io/client-go/util/workqueue"
	"time"
	"k8s.io/apimachinery/pkg/util/wait"
	"github.com/golang/glog"
	"k8s.io/client-go/tools/cache"
	"fmt"
)

var (
	keyFunc = cache.DeletionHandlingMetaNamespaceKeyFunc
)

type Queue struct {
	queue workqueue.RateLimitingInterface

	sync func(interface{}) error

	workerDone chan bool

	fn func(obj interface{}) (interface{}, error)

	lastSync int64
}

type Element struct{
	Key 		interface{}
	Timestamp 	int64
}

func (t *Queue) Run(period time.Duration, stopCh <-chan struct{})  {
	wait.Until(t.worker, period, stopCh)
}

func (t *Queue) Enqueue(obj interface{})  {
	if t.IsShuttingDown() {
		glog.Errorf("Queue has been shutdown, failed to enqueue: %v", obj)
	}

	ts := time.Now().UnixNano()
	glog.V(3).Infof("queuing item %v", obj)
	key, err := t.fn(obj)
	if err != nil {
		glog.Errorf("%v", err)
		return
	}
	t.queue.Add(Element{
		Key:       key,
		Timestamp: ts,
	})

}

func (t *Queue) worker()  {
	for {
		key, quit := t.queue.Get()
		if quit {
			if !isClosed(t.workerDone) {
				close(t.workerDone)
			}
			return
		}
		ts := time.Now().UnixNano()

		item := key.(Element)
		if t.lastSync > item.Timestamp {
			glog.V(3).Infof("skipping %v sync (%v > %v)", item.Key, t.lastSync, item.Timestamp)
			t.queue.Forget(key)
			t.queue.Done(key)
			continue
		}

		glog.V(3).Infof("syncing %v", item.Key)
		err := t.sync(key)
		if err != nil {
			glog.Warning("requeuing %v because of %v", item.Key, err)
			t.queue.AddRateLimited(Element{
				Key:			item.Key,
				Timestamp:		time.Now().UnixNano(),
			})
		} else {
			t.queue.Forget(key)
			t.lastSync = ts
		}

		t.queue.Done(key)
	}
}

func (t *Queue) defaultKeyFunc(obj interface{}) (interface{}, error) {
	key, err := keyFunc(obj)
	if err != nil {
		return "", fmt.Errorf("couldn't get key from object %+v: %v", obj, err)
	}

	return key, nil
}

func (t *Queue) ShutDown()  {
	t.queue.ShutDown()
	<-t.workerDone
}

func (t *Queue) IsShuttingDown() bool {
	return t.queue.ShuttingDown()
}

func isClosed(ch <-chan bool) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func NewTaskQueue(syncFn func(interface{}) error) *Queue {
	return NewCustomTaskQueue(syncFn, nil)
}

func NewCustomTaskQueue(syncFn func(interface{}) error, fn func(interface{}) (interface{}, error)) *Queue {
	q := &Queue{
		queue: 			workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		sync:			syncFn,
		workerDone:		make(chan bool),
		fn:				fn,
	}

	if fn == nil {
		q.fn = q.defaultKeyFunc
	}

	return q
}