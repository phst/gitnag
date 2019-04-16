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
	"os/exec"
	"strings"

	"github.com/golang/glog"
)

type repo struct {
	workTree string
	bare     bool
}

func getRepo(dir string) (*repo, error) {
	workTree, err := getWorkTree(dir)
	if err != nil {
		return nil, err
	}
	if workTree == "" {
		return nil, nil
	}
	bare, err := isBare(dir)
	if err != nil {
		return nil, err
	}
	r := &repo{
		workTree: workTree,
		bare:     bare,
	}
	return r, nil
}

func isBare(dir string) (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-bare-repository")
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("error running %v: %v", cmd.Args, err)
	}
	switch s := string(out); s {
	case "true\n":
		return true, nil
	case "false\n":
		return false, nil
	default:
		return false, fmt.Errorf("unexpected output of %v: %q", cmd.Args, s)
	}
}

func getWorkTree(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.Exited() && !exitErr.Success() {
				return "", nil
			}
		}
		return "", fmt.Errorf("error running %v: %v", cmd.Args, err)
	}
	path := strings.TrimSuffix(string(out), "\n")
	if path == "" {
		return "", fmt.Errorf("%v returned empty output", cmd.Args)
	}
	glog.Infof("found work tree in %s", path)
	return path, nil
}

func (r repo) call(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.workTree
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running %v: %v", cmd.Args, err)
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}
