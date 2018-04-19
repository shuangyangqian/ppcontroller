package main

import (
	"podtoservice/pkg/controller"
	"github.com/spf13/pflag"
	"flag"
	"os"
)

func parseFlags() (*controller.Configuration, error)  {
	var (
		flags = pflag.NewFlagSet("", pflag.ExitOnError)

		kubeConfigFile = flags.String("kubeconfig", "", "Path of config file to connect the k8s apiServer")
	)

	flag.Set("logtostderr", "true")
	flags.AddGoFlagSet(flag.CommandLine)
	flags.Parse(os.Args)
	//flag.Set("logtostderr", "true")

	//flag.CommandLine.Parse([]string{})

	config := &controller.Configuration{
		kubeConfigFile:			*kubeConfigFile,
	}

	return config, nil
}