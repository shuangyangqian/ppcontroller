package main

import (
	"flag"
	"os"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	p, err := os.Getwd()

	if err != nil {
		glog.Info("Getwd: ", err)
	} else {
		glog.Info("Getwd: ", p)
	}

	glog.Info("Prepare to repel boarders")
	glog.Info("222222222222---log_backtrace_at")
	glog.Info("333333333333")

	glog.V(1).Infoln("Processed1", "nItems1", "elements1")

	glog.V(2).Infoln("Processed2", "nItems2", "elements2")

	glog.V(3).Infoln("Processed3", "nItems3", "elements3")

	glog.V(4).Infoln("Processed4", "nItems4", "elements4")

	glog.V(5).Infoln("Processed5", "nItems5", "elements5")

	glog.Error("errrrr")


	exit()
}

func exit()  {
	glog.Flush()
}