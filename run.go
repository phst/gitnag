// Copyright 2015, 2022 Google LLC
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
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
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
			log.Printf("found problem: %s", p)
			problems = append(problems, p)
		}
		return filepath.SkipDir
	}
	for dir := range cfg.Directories {
		log.Printf("processing directory %s", dir)
		if err := filepath.Walk(dir, runDir); err != nil {
			return fmt.Errorf("could not walk directory %s: %v", dir, err)
		}
	}
	log.Printf("found %d problems", len(problems))
	sort.Sort(problems)
	for _, p := range problems {
		if err := p.notify(); err != nil {
			return err
		}
	}
	log.Print("completed successfully")
	return nil
}

type problem struct {
	repo
	desc string
}

func (p problem) String() string {
	return fmt.Sprintf("%s: %s", shortenFilename(p.workTree), p.desc)
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

func shortenFilename(s string) string {
	s = filepath.Clean(s)
	u, err := user.Current()
	if err != nil || u.HomeDir == "" {
		return s
	}
	h := filepath.Clean(u.HomeDir)
	p := h + string(filepath.Separator)
	if strings.HasPrefix(s, p) {
		return filepath.Join("~", s[len(p):])
	}
	return s
}
