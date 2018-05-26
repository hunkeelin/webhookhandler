package main

import (
    "github.com/hunkeelin/klinenv"
    "path/filepath"
    "errors"
    "strings"
)

func Determine(p string, g Gitpayload,d string) (string,[]JobConfig, error) {
    var secret string
    var out []JobConfig
    var fs string
    if d == "repo" {
        fs = g.Repository.URL
    } else {
        fs = g.Repository.Owner.HTMLURL
    }
    if fs == ""{
        return secret,out,errors.New("Repo URL doesn't exist in payload")
    }
    path := p + strings.TrimPrefix(fs, "https://")
    check,_ := exist(path+"/"+"config")
    if !check{
        return secret,out,errors.New("Config and path doesn't exist \n"+"P:"+path+"\n"+"F:"+path+"/config")
    }
    // return secret
    jobconfig := klinenv.NewAppConfig(path+"/"+"config")
    s,err := jobconfig.Get("secret")
    if err != nil {
        return secret,out, err
    }
    secret = s
    // return slice struct
    jobpaths, err := filepath.Glob(path + "/" + "*.conf")
    if err != nil {
        return secret,out,err
    }
    for _,jobpath := range jobpaths {
        job := klinenv.NewAppConfig(jobpath)
        run, err := job.Get("run")
        if err != nil {
            return secret,out, err
        }
        j := JobConfig{run: run}
        out = append(out,j)
    }
    return secret,out, nil
}
