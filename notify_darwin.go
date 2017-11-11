package main

import (
	"os/exec"

	shellquote "github.com/kballard/go-shellquote"
)

func (p problem) notify() error {
	cmd := exec.Command(
		"terminal-notifier", "-message", p.String(), "-title", "gitnag",
		"-execute", shellquote.Join("open", "-b", "com.apple.Terminal", "--", p.workTree))
	return cmd.Run()
}
