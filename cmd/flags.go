package main

import (
	"ppcontroller/pkg/controller"
	"github.com/spf13/pflag"
	"flag"
	"os"
)

func parseFlags() (*controller.Configuration, error)  {
	var (
		flags = pflag.NewFlagSet("", pflag.ExitOnError)

		kubeConfigFile = flags.String("kubeconfig", "", "Path of config file to connect the k8s apiServer")
	)

	//flag.Set("logtostderr", "true")

	// change the log file dir to /tmp
	flag.Set("log_dir", "/tmp")
	flags.AddGoFlagSet(flag.CommandLine)
	flags.Parse(os.Args)

	//flag.CommandLine.Parse([]string{})

	config := &controller.Configuration{
		KubeConfigFile:			*kubeConfigFile,
	}

	return config, nil
}