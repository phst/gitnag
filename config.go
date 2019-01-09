package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/golang/glog"
)

type config struct {
	Directories map[string]struct{}
}

func loadConfig(file string) (*config, error) {
	glog.Infof("loading configuration from %v", file)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var r config
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	glog.Infof("configuration specifies %d directories", len(r.Directories))
	return &r, nil
}
