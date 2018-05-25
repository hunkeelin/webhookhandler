package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
    switch apikey := r.Header.Get("api-key"); apikey {
        case "": //github json by defualt
            msg,status := GitWork(r,f)
            w.WriteHeader(status)
            w.Write([]byte(msg))
            fmt.Println(msg,status)
            return
        case f.apikey: // application specific.
            fmt.Println("right apikey")
            return
        default: // application specific but wrong api key. 
            fmt.Println("wrong apikey")
            w.WriteHeader(300)
            w.Write([]byte("Wrong apikey\n"))
            return
    }
    //f.mu.Lock()
    //dowork()
    ////f.mu.Unlock()
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
	//err := s.ListenAndServeTLS(c.certpath, c.keypath)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("can't listen and serve check port and binding addr")
	}
}
