package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func CheckSecret(rs string, r *http.Request, body []byte) (error, string, int) {
	secret := []byte(rs)
	if !verifySignature(secret, r.Header.Get("X-Hub-Signature"), body) {
		return errors.New("non nil"), "Bad Signiture is the secret correct in the config?", 402
	}
	if !isvalidmethod(r) {
		return errors.New("non nil"), "Bad request method " + r.Method, 405
	}
	return nil, "", 0
}
func GitExec(t []JobConfig, f *Conn) {
	for _, task := range t {
		cmd := "sh"
		args := []string{task.run}
		err := runshell(cmd, args, f.uid, f.gid)
		if err != nil {
			fmt.Println(args)
			fmt.Println(err)
		}
	}
	//    f.sem <- struct{}{}
	//    dowork()
	//    <-f.sem
}

func GitWork(r *http.Request, f *Conn) (string, int) {
	if r.Header.Get("content-type") != "application/json" {
		return "Bad Request,the payload is not application/json ", 400
	}
	limitbody := io.LimitReader(r.Body, 65535)
	body, _ := ioutil.ReadAll(limitbody)
	var g Gitpayload
	b := bytes.NewReader(body)
	err := json.NewDecoder(b).Decode(&g)
	if err != nil {
		return "Bad Request,the payload is not a valid json ", 400
	}
	rs, t, err := Determine(f.jobdir, g, "org")
	if err != nil {
		fmt.Println(err)
		//    return "Error reading config file: ",300
	} else {
		err, st, in := CheckSecret(rs, r, body)
		if err != nil {
			return st, in
		}
		GitExec(t, f)
	}
	rs, t, err = Determine(f.jobdir, g, "repo")
	if err != nil {
		fmt.Println(err)
		return "Error reading config file: ", 300
	} else {
		err, st, in := CheckSecret(rs, r, body)
		if err != nil {
			return st, in
		}
		GitExec(t, f)
	}
	return "Status -ok ", 200
}
