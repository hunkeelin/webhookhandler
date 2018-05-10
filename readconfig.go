package main

import (
    "github.com/hunkeelin/klinenv"
    "log"
    "strconv"
)

func readconfig() (toreturn []string, concur int) {
    var to_return []string

    config := klinenv.NewAppConfig("genkins.conf")

    certdir, err := config.Get("certpath")
    if err != nil {
        log.Fatal("unable to retrieve the value of certdir check config file /etc/genkins/genkins.conf")
    }
    to_return = append(to_return,certdir)
    keydir, err := config.Get("keypath")
    if err != nil {
        log.Fatal("unable to retrieve the value of certdir check config file /etc/genkins/genkins.conf")
    }
    to_return = append(to_return,keydir)

    remote, err := config.Get("giturl")
    if err != nil {
        log.Fatal("unable to retrieve the value of remote check config file /etc/genkins/genkins.conf")
    }
    to_return = append(to_return,remote)

    lockfile, err := config.Get("lockfile")
    if err != nil {
        log.Fatal("unable to retrieve the value of lockfile check config file /etc/genkins/genkins.conf")
    }
    to_return = append(to_return,lockfile)

    bindaddr, err := config.Get("bindaddr")
    if err != nil {
        log.Fatal("unable to retrieve the value of bindaddr check config file /etc/genkins/genkins.conf")
    }
    to_return = append(to_return,bindaddr)

    port, err := config.Get("port")
    if err != nil {
        log.Fatal("unable to retrieve the value of bindaddr check config file /etc/genkins/genkins.conf")
    }
    to_return = append(to_return,port)

    rconcur, err := config.Get("concur")
    if err != nil {
        log.Fatal("unable to retrieve the value of concur check config file /etc/genkins/genkins.conf")
    }
    to_return_int, err := strconv.Atoi(rconcur)
    if err != nil {
        log.Fatal("can't convert string to int for concur in /etc/genkins/genkins.conf")
    }
    return to_return, to_return_int
}
