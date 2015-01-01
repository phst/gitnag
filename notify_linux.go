package main

import "os/exec"

func (p problem) notify() error {
	cmd := exec.Command("notify-send", "--", "gitnag", p.String())
	return cmd.Run()
}
