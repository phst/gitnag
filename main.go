package main

import (
	"flag"

	"github.com/golang/glog"
)

var configFile = flag.String("config", "", "location of the configuration file")

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		glog.Fatal("this program doesn't accept positional arguments")
	}
	if *configFile == "" {
		glog.Fatal("-config flag must be provided")
	}
	cfg, err := loadConfig(*configFile)
	if err != nil {
		glog.Fatalf("could not load configuration: %v", err)
	}
	if err := run(*cfg); err != nil {
		glog.Fatal(err)
	}
}
