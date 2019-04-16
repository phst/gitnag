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
	"os"
	"os/exec"
	"strings"
)

func (p problem) notify() error {
	cmd := exec.Command("notify-send", "--", "gitnag", p.String())
	cmd.Env = addDisplayEnv(os.Environ())
	return cmd.Run()
}

func addDisplayEnv(env []string) []string {
	for _, s := range env {
		if strings.HasPrefix(s, "DISPLAY=") {
			return env
		}
	}
	return append(env, "DISPLAY=:0")
}
