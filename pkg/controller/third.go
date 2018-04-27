package controller

import (
	"fmt"
	"ppcontroller/pkg/task"
	"k8s.io/api/core/v1"
	"net/http"
	"io/ioutil"
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

func (pts *PODController) httpGet()  {
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		fmt.Printf("We get the err: %v\n", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("We get the err: %v\n", err)
	}

	fmt.Println(string(body))
}