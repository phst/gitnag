// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
