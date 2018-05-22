package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
	limitbody := io.LimitReader(r.Body, 65535)
	body, _ := ioutil.ReadAll(limitbody)
	var g Gitpayload
	b := bytes.NewReader(body)
	err := json.NewDecoder(b).Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request,the payload is not a valid json\n"))
		return
	}
	rs, t, err := Determine(f.jobdir, g)
	if err != nil {
		fmt.Println("error reading config file:", err)
		return
	}
	secret := []byte(rs)
	if !verifySignature(secret, r.Header.Get("X-Hub-Signature"), body) {
		w.WriteHeader(402)
		w.Write([]byte("Bad signiture \n"))
		return
	}
	if !isvalidmethod(r) {
		w.WriteHeader(405)
		w.Write([]byte("Bad request method " + r.Method + "\n"))
		return
	}
	// logic
    for _,task := range t {
        cmd := "sh"
        args := []string{task.run}
        err := runshell(cmd,args)
        if err != nil{
            fmt.Println(task.run)
            fmt.Println(err)
        }
    }
    //f.mu.Lock()
    f.sem <- struct{}{}
    dowork()
    <-f.sem
    ////f.mu.Unlock()
	w.WriteHeader(200)
	w.Write([]byte("Status ok: do logic\n"))
	return
}
func main() {
	newcon := new(Conn)
	// define config params
	c := readconfig()
    sema := make(chan struct {},1)
    newcon.sem = sema
	newcon.apikey = c.apikey
	newcon.concur = c.concur
	newcon.jobdir = c.jobdir

	tlsconfig := &tls.Config{
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
	}
	con := http.NewServeMux()
	con.HandleFunc("/", newcon.handleWebHook)
	s := &http.Server{
		Addr:         c.bindaddr + ":" + c.port,
		TLSConfig:    tlsconfig,
		Handler:      con,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	fmt.Println("listening to " + c.bindaddr + " " + c.port)
	err := s.ListenAndServeTLS(c.certpath, c.keypath)
	//err := s.ListenAndServe()
	if err != nil {
		log.Fatal("can't listen and serve check port and binding addr")
	}
}
