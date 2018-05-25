package main

import (
	"github.com/hunkeelin/klinenv"
	"log"
	"strconv"
	"strings"
)

func readconfig() Config {
	config := klinenv.NewAppConfig("/etc/genkins/genkins.conf")
	rconcur, err := config.Get("concur")
	if err != nil {
		log.Fatal("unable to retrieve the value of concur check config file")
	}
	concur, err := strconv.Atoi(rconcur)
	if err != nil {
		log.Fatal("can't convert string to int for concur")
	}
	var c Config
	apikey, err := config.Get("apikey")
	checkerr(err)
	c.apikey = apikey

	bindaddr, err := config.Get("bindaddr")
	checkerr(err)
	c.bindaddr = bindaddr

	port, err := config.Get("port")
	checkerr(err)
	c.port = port

	certpath, err := config.Get("certpath")
	checkerr(err)
	c.certpath = certpath

	keypath, err := config.Get("keypath")
	checkerr(err)
	c.keypath = keypath

	hosts, err := config.Get("hosts")
	checkerr(err)
	c.hosts = strings.Split(hosts, ",")

	jobdir, err := config.Get("jobdir")
	checkerr(err)
	if string(jobdir[len(jobdir)-1]) != "/" {
		jobdir = jobdir + "/"
	}
	c.jobdir = jobdir

	c.concur = concur
	return c
}
