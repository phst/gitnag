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
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/golang/glog"
)

func run(cfg config) error {
	var problems problems
	runDir := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		r, err := getRepo(path)
		if err != nil {
			return err
		}
		if r == nil {
			return nil
		}
		if err := r.check(); err != nil {
			p := problem{*r, err.Error()}
			glog.Errorf("found problem: %s", p)
			problems = append(problems, p)
		}
		return filepath.SkipDir
	}
	for dir := range cfg.Directories {
		glog.Infof("processing directory %s", dir)
		if err := filepath.Walk(dir, runDir); err != nil {
			return fmt.Errorf("could not walk directory %s: %v", dir, err)
		}
	}
	sort.Sort(problems)
	for _, p := range problems {
		if err := p.notify(); err != nil {
			return err
		}
	}
	return nil
}

type problem struct {
	repo
	desc string
}

func (p problem) String() string {
	return fmt.Sprintf("%s: %s", p.workTree, p.desc)
}

type problems []problem

func (p problems) Len() int {
	return len(p)
}

func (p problems) Less(i, j int) bool {
	if p[i].workTree != p[j].workTree {
		return p[i].workTree < p[j].workTree
	}
	return p[i].desc < p[j].desc
}

func (p problems) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
