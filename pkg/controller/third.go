package controller

import (
	"fmt"
	"ppcontroller/pkg/task"
	"k8s.io/api/core/v1"
)

func (pts *PODController) printMsg(item interface{}) error {
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