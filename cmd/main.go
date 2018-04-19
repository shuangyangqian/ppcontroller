package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"podtoservice/pkg/controller"
	"github.com/golang/glog"
)

func main() {
	conf, err := parseFlags()

	if err != nil {
		glog.Errorf("we got some unexpect err: %v", err)
	}
	clientset := getKubeClientset(conf.KubeConfigFile)
	conf.Client = clientset

	controller := controller.NewPodController(conf)
	var stopCh <-chan struct{}
	controller.Run(2, stopCh)

}

func getKubeClientset(kubeconfig string) kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return *clientset
}
