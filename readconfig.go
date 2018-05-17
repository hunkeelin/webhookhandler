package main

import (
	"github.com/hunkeelin/klinenv"
	"log"
	"strconv"
)

func readconfig() (toreturn klinenv.AppConfig, concur int) {
	config := klinenv.NewAppConfig("genkins.conf")
	rconcur, err := config.Get("concur")
	if err != nil {
		log.Fatal("unable to retrieve the value of concur check config file")
	}
	to_return_int, err := strconv.Atoi(rconcur)
	if err != nil {
		log.Fatal("can't convert string to int for concur")
	}
	return config, to_return_int
}
