package main

import (
	"github.com/hunkeelin/klinenv"
	"path/filepath"
	"strings"
)

type Jobconfig struct {
	secret string
	cmd    string
}

func Determine(p string, g Gitpayload) (*Jobconfig, error) {
	j := new(Jobconfig)
	path := p + strings.TrimPrefix(g.Repository.URL, "https://")
	jobpath, _ := filepath.Glob(path + "/" + "*")
	job := klinenv.NewAppConfig(jobpath[0])
	secret, err := job.Get("secret")
	if err != nil {
		return j, err
	}
	cmd, err := job.Get("cmd")
	if err != nil {
		return j, err
	}
	j.secret = secret
	j.cmd = cmd
	return j, nil
}
