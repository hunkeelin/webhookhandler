package main

import (
	"github.com/hunkeelin/klinenv"
	"log"
	"strconv"
)

func readconfig() (toreturn map[string]string, concur int) {
	to_return := make(map[string]string)
	confvalues := []string{"lockfile", "giturl", "certpath", "keypath", "bindaddr", "port"}
	config := klinenv.NewAppConfig("genkins.conf")

	for _, element := range confvalues {
		confval, err := config.Get(element)
		if err != nil {
			log.Fatal("unable to retrieve the value of " + element + "check config file ")
		}
		to_return[element] = confval
	}

	rconcur, err := config.Get("concur")
	if err != nil {
		log.Fatal("unable to retrieve the value of concur check config file")
	}
	to_return_int, err := strconv.Atoi(rconcur)
	if err != nil {
		log.Fatal("can't convert string to int for concur")
	}
	return to_return, to_return_int
}
