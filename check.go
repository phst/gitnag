package main

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

func (r repo) check() error {
	if !r.bare {
		err := r.workTreeClean()
		if err != nil {
			return err
		}
	}
	if err := r.pull(); err != nil {
		return err
	}
	return r.allCommitsPushed()
}

func (r repo) workTreeClean() error {
	glog.Infof("checking whether work tree and index in %s are clean", r.workTree)
	out, err := r.call("status", "-z")
	if err != nil {
		return err
	}
	if out != "" {
		return fmt.Errorf("unclean work tree or index in %s", r.workTree)
	}
	return nil
}

func (r repo) pull() error {
	glog.Infof("pulling to %s", r.workTree)
	if _, err := r.call("pull", "--ff-only"); err != nil {
		return fmt.Errorf("could not pull %s: %v", r.workTree, err)
	}
	return nil

}

func (r repo) allCommitsPushed() error {
	glog.Infof("comparing commits of %s to upstream", r.workTree)
	out, err := r.call("rev-list", "--count", "@{upstream}..")
	if err != nil {
		return err
	}
	n, err := strconv.Atoi(out)
	if err != nil {
		return fmt.Errorf("could not parse commit count %q: %v", out, err)
	}
	if n != 0 {
		return fmt.Errorf("%d commits unpushed in %s", n, r.workTree)
	}
	return nil
}
