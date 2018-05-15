package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//	"time"
)

func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
	if !isvalidmethod(r) {
		w.WriteHeader(405)
		w.Write([]byte("Bad request method " + r.Method + "\n"))
		return
	}
	var g Gitpayload
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request,the payload is not a valid json\n"))
		return
	}

	fmt.Println(r.Header.Get("api-key"))
	if r.Header.Get("api-key") != "password" {
		fmt.Println("wrong api key")
		return
	}
	if g.doesmatchbody(f.regex) {
		w.WriteHeader(200)
		w.Write([]byte("Matches\n"))
		return
	} else {
		w.WriteHeader(400)
		w.Write([]byte("Doesn't match " + f.regex))
		return
	}
	return
}
func main() {
	newcon := new(Conn)
	config, _ := readconfig()
	newcon.regex = config["giturl"]
	http.HandleFunc("/", newcon.handleWebHook)
	//	s := &http.Server{
	//		Addr:         ":8080",
	//		ReadTimeout:  5 * time.Second,
	//		WriteTimeout: 10 * time.Second,
	//		IdleTimeout:  120 * time.Second,
	//	}
	//	log.Fatal(s.ListenAndServeTLS("", ""))
	err := http.ListenAndServeTLS(config["bindaddr"]+":"+config["port"], config["certpath"], config["keypath"], nil)
	if err != nil {
		log.Fatal("Unable to Listen to port", err)
	}
}
