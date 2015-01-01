package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	shellquote "github.com/kballard/go-shellquote"
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
			problems = append(problems, problem{*r, err.Error()})
		}
		return filepath.SkipDir
	}
	for dir := range cfg.Directories {
		log.Printf("processing directory %s", dir)
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

func (p problem) dir() string {
	if p.workTree != "" {
		return p.workTree
	}
	return p.gitDir
}

func (p problem) String() string {
	return fmt.Sprintf("%s: %s", p.dir(), p.desc)
}

type problems []problem

func (p problems) Len() int {
	return len(p)
}

func (p problems) Less(i, j int) bool {
	if p[i].gitDir != p[j].gitDir {
		return p[i].gitDir < p[j].gitDir
	}
	return p[i].desc < p[j].desc
}

func (p problems) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p problem) notify() error {
	cmd := exec.Command(
		"terminal-notifier", "-message", p.String(), "-title", "gitnag",
		"-execute", shellquote.Join("open", "-b", "com.apple.Terminal", "--", p.dir()))
	return cmd.Run()
}
