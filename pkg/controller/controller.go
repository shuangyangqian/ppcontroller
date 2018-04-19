package controller

import (
	"k8s.io/client-go/kubernetes"
	"time"
	"fmt"
	"k8s.io/apimachinery/pkg/util/runtime"
	"podtoservice/pkg/task"
	"k8s.io/api/core/v1"
	"github.com/golang/glog"
)

type podController struct {

	kubeClient kubernetes.Clientset

	stopChan chan struct{}

	podListener *cacheLister
	podController *cacheController

	syncQueue *task.Queue

}

func NewPodController(conf *Configuration) *podController {
	pts := &podController{
		kubeClient: conf.Client,
		stopChan:make(chan struct{}),
	}

	pts.podListener, pts.podController = pts.createListers(pts.stopChan)

	pts.syncQueue = task.NewTaskQueue(pts.printMsg)

	return pts
}

func (pts *podController) Run(workers int, stpChan <- chan struct{})  {
	defer runtime.HandleCrash()
	glog.Info("Starting podToService controlelr Manager...")
	pts.podController.Run(pts.stopChan)

	go pts.syncQueue.Run(time.Second, pts.stopChan)

	for {
		select {
		case <-pts.stopChan:
			glog.Info("stop is received")
		default:
		}
	}

}

func (pts *podController) printMsg(item interface{}) error {
	if pts.syncQueue.IsShuttingDown() {
		return nil
	}

	if element, ok := item.(task.Element); ok {
		if name, ok := element.Key.(string); ok {
			if obj, exists, _ := pts.podListener.pod.GetByKey(name); exists {
				pod := obj.(*v1.Pod)
				fmt.Printf("we get the change of POD:%v\n", pod)
			}
		}
	}
	return nil
}


type Configuration struct {
	KubeConfigFile string
	Client kubernetes.Clientset
}