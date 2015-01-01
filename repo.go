package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type repo struct {
	gitDir, workTree string
	bare             bool
}

func getRepo(dir string) (*repo, error) {
	gitDir, err := getGitDir(dir)
	if err != nil {
		return nil, err
	}
	if gitDir == "" {
		return nil, nil
	}
	bare, err := isBare(gitDir)
	if err != nil {
		return nil, err
	}
	workTree, err := getWorkTree(dir)
	if err != nil {
		return nil, err
	}
	r := &repo{
		gitDir:   gitDir,
		workTree: workTree,
		bare:     bare,
	}
	return r, nil
}

func getGitDir(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
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
	gitDir := strings.TrimSuffix(string(out), "\n")
	if gitDir == "" {
		return "", fmt.Errorf("%v returned empty output", cmd.Args)
	}
	gitDir = filepath.Join(dir, gitDir)
	log.Printf("found Git repository in %s", gitDir)
	return gitDir, nil
}

func isBare(dir string) (bool, error) {
	cmd := exec.Command("git", "--git-dir", dir, "rev-parse", "--is-bare-repository")
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
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
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
	args = append([]string{"--git-dir", r.gitDir, "--work-tree", r.workTree}, args...)
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running %v: %v", cmd.Args, err)
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}
