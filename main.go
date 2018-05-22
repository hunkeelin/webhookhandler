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
	"sync"
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
	d, err := Determine(f.jobdir, g)
	if err != nil {
		fmt.Println("error reading config file ", err)
		return
	}

	secret := []byte(d.secret)
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
	sema := make(chan struct{}, f.concur)
	wg := sync.WaitGroup{}
	for _, i := range f.hosts {
		sema <- struct{}{}
		wg.Add(1)
		go func(g string) {
			<-sema
			dowork(g)
			wg.Done()
		}(i)
	}
	fmt.Println("dowork")
	w.WriteHeader(200)
	w.Write([]byte("Status ok: do logic\n"))
	return
}
func main() {
	newcon := new(Conn)
	// define config params
	c := readconfig()
	newcon.apikey = c.apikey
	newcon.concur = c.concur
	newcon.hosts = c.hosts
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
	//	err := s.ListenAndServeTLS(c.certpath, c.keypath)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("can't listen and serve check port and binding addr")
	}
}
