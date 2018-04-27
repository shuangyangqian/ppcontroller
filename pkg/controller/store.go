package controller

import (
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/util/runtime"
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
)

type cacheController struct {
	Pod cache.Controller
}

type cacheLister struct {
	pod cache.Store
}

func (c *cacheController) Run(stopCh chan struct{}) {
	go c.Pod.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh,
		c.Pod.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timeout waiting for caches to sync"))
	}
}

func (pc *PODController) createListers(stopChan chan struct{}) (*cacheLister, *cacheController) {
	podEventHander := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("some obj:%v has been added to the cluster", obj)
			pc.syncQueue.Enqueue(obj)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("some obj:%v has been deleted from the cluster", obj)
			pc.syncQueue.Enqueue(obj)
		},
	}

	watchNs := v1.NamespaceAll

	lister := &cacheLister{}
	controller := &cacheController{}
	lister.pod, controller.Pod = cache.NewInformer(
		cache.NewListWatchFromClient(pc.kubeClient.CoreV1().RESTClient(), "pods", watchNs, fields.Everything()),
		&v1.Pod{}, 0, podEventHander)

	return lister, controller
}
