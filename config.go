package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type config struct {
	Directories map[string]struct{}
}

func loadConfig(file string) (*config, error) {
	log.Printf("loading configuration from %v", file)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var r config
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	log.Printf("configuration specifies %d directories", len(r.Directories))
	return &r, nil
}
