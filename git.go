package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func GitWork(r *http.Request,f *Conn)(string,int) {
    if r.Header.Get("content-type") != "application/json"{
        return "Bad Request,the payload is not application/json ",400
    }
    limitbody := io.LimitReader(r.Body, 65535)
    body, _ := ioutil.ReadAll(limitbody)
    var g Gitpayload
    b := bytes.NewReader(body)
    err := json.NewDecoder(b).Decode(&g)
    if err != nil {
        return "Bad Request,the payload is not a valid json ",400
    }
    rs, t, err := Determine(f.jobdir, g)
    if err != nil {
        fmt.Println(err)
        return "Error reading config file: ",300
    }
    secret := []byte(rs)
    if !verifySignature(secret, r.Header.Get("X-Hub-Signature"), body) {
        return "Bad Signiture",402
    }
    if !isvalidmethod(r) {
        return "Bad request method " + r.Method,405
    }
    for _,task := range t {
        cmd := "sh"
        args := []string{task.run}
        err := runshell(cmd,args)
        if err != nil{
            fmt.Println(task.run)
            fmt.Println(err)
        }
    }
    f.sem <- struct{}{}
    dowork()
    <-f.sem
    return "Status -ok ", 200
}
