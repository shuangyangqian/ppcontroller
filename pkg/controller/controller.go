package controller

import (
	"k8s.io/client-go/kubernetes"
	"time"
	"k8s.io/apimachinery/pkg/util/runtime"
	"ppcontroller/pkg/task"
	"github.com/golang/glog"
)

type PODController struct {

	kubeClient kubernetes.Clientset

	stopChan chan struct{}

	podListener *cacheLister
	podController *cacheController

	syncQueue *task.Queue

}

func NewPodController(conf *Configuration) *PODController {
	pts := &PODController{
		kubeClient: conf.Client,
		stopChan:make(chan struct{}),
	}

	pts.podListener, pts.podController = pts.createListers(pts.stopChan)

	pts.syncQueue = task.NewTaskQueue(pts.printMsg)

	return pts
}

func (pts *PODController) Start()  {
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

type Configuration struct {
	KubeConfigFile string
	Client kubernetes.Clientset
}