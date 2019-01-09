package main

import (
	"flag"

	"github.com/golang/glog"
)

var configFile = flag.String("config", "", "location of the configuration file")

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		glog.Exit("this program doesn't accept positional arguments")
	}
	if *configFile == "" {
		glog.Exit("-config flag must be provided")
	}
	cfg, err := loadConfig(*configFile)
	if err != nil {
		glog.Exitf("could not load configuration: %v", err)
	}
	if err := run(*cfg); err != nil {
		glog.Exit(err)
	}
}
