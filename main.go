package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
	var gitpayload Gitpayload
	err := json.NewDecoder(r.Body).Decode(&gitpayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request,the payload json is not correct please read the documentation"))
	}
	if gitpayload.doesmatchbody(f.regex) {
		fmt.Println("match")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Doesn't match " + f.regex))
	}
	fmt.Println(r.Header)
}
func main() {
	newcon := new(Conn)
	config, _ := readconfig()
	newcon.regex = config["giturl"]
	http.HandleFunc("/", newcon.handleWebHook)
	err := http.ListenAndServeTLS(config["bindaddr"]+":"+config["port"], config["certpath"], config["keypath"], nil)
	if err != nil {
		log.Fatal("Unable to Listen to port", err)
	}
}
