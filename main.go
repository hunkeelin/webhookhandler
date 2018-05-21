package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
    "bytes"
    "io"
    "io/ioutil"
	"net/http"
	"time"
)

func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
    secret := []byte(f.secret)
    limitbody := io.LimitReader(r.Body,65535)
    body, _ := ioutil.ReadAll(limitbody)
    if !verifySignature(secret,r.Header.Get("X-Hub-Signature"),body){
        w.WriteHeader(402)
		w.Write([]byte("Bad signiture \n"))
		return
    }
	if !isvalidmethod(r) {
		w.WriteHeader(405)
		w.Write([]byte("Bad request method " + r.Method + "\n"))
		return
	}
	var g Gitpayload
    b := bytes.NewReader(body)
	err := json.NewDecoder(b).Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request,the payload is not a valid json\n"))
		return
	}
//	if r.Header.Get("api-key") != f.apikey {
//		w.WriteHeader(405)
//		w.Write([]byte("Wrong apikey\n"))
//		return
//	}
	if !g.doesmatchbody(f.regex) {
		w.WriteHeader(400)
		w.Write([]byte("Doesn't Match " + f.regex))
		return
	}
	// logic
	w.WriteHeader(200)
	w.Write([]byte("Status ok: do logic\n"))
	return
}
func main() {
	newcon := new(Conn)
    // define config params
	config, _ := readconfig()
	regex , err  := config.Get("giturl")
    checkerr(err)
    newcon.regex = regex
	apikey, err := config.Get("apikey")
    checkerr(err)
    newcon.apikey = apikey
    binddir , err := config.Get("bindaddr")
    checkerr(err)
    addrport, err := config.Get("port")
    checkerr(err)
    certpath, err := config.Get("certpath")
    checkerr(err)
    keypath, err := config.Get("keypath")
    checkerr(err)
    secret, err := config.Get("secret")
    checkerr(err)
    newcon.secret = secret
    // end of define config params
	tlsconfig := &tls.Config{
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},
	}
	c := http.NewServeMux()
	c.HandleFunc("/", newcon.handleWebHook)
	s := &http.Server{
		Addr:         binddir + ":" + addrport,
		TLSConfig:    tlsconfig,
		Handler:      c,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
    fmt.Println("listening to " + binddir + " " + addrport)
	err = s.ListenAndServeTLS(certpath, keypath)
	if err != nil {
		log.Fatal("can't listen and serve check port and binding addr")
	}
}
