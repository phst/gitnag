package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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
	log.Printf("found work tree in %s", path)
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
