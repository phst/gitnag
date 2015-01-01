package main

import (
	"flag"
	"log"
)

var configFile = flag.String("config", "", "location of the configuration file")

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		log.Fatal("this program doesn't accept positional arguments")
	}
	if *configFile == "" {
		log.Fatal("-config flag must be provided")
	}
	cfg, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("could not load configuration: %v", err)
	}
	if err := run(*cfg); err != nil {
		log.Fatal(err)
	}
}
