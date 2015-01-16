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
