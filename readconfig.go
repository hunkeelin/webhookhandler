package main

import (
	"github.com/hunkeelin/klinenv"
	"log"
	"os/user"
	"strconv"
	"strings"
)

func readconfig(p string) Config {
	config := klinenv.NewAppConfig(p)
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

	rawuser, err := config.Get("user")
	checkerr(err)
	if rawuser == "root" || rawuser == "" {
		log.Fatal("invalid user, please specify a non root user in genkins.conf")
	}
	userdata, err := user.Lookup(rawuser)
	checkerr(err)

	uid, err := strconv.ParseUint(userdata.Uid, 10, 32)
	checkerr(err)

	c.uid = uint32(uid)
	gid, err := strconv.ParseUint(userdata.Gid, 10, 32)

	checkerr(err)
	c.gid = uint32(gid)

	c.homedir = userdata.HomeDir

	hosts, err := config.Get("hosts")
	checkerr(err)
	c.hosts = strings.Split(hosts, ",")

	jobdir, err := config.Get("jobdir")
	checkerr(err)
    if len(jobdir) == 0 {
        c.jobdir = jobdir
    } else {
        if string(jobdir[len(jobdir)-1]) != "/" {
            jobdir = jobdir + "/"
        }
        c.jobdir = jobdir
    }

	c.concur = concur
	return c
}
