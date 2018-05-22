package main

import (
    "github.com/hunkeelin/klinenv"
    "path/filepath"
    "errors"
    "strings"
)

func Determine(p string, g Gitpayload) (string,[]JobConfig, error) {
    var secret string
    var out []JobConfig
    path := p + strings.TrimPrefix(g.Repository.URL, "https://")
    check,_ := exist(path)
    if !check{
        return secret,out,errors.New("path doesn't exist")
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
        cmd, err := job.Get("cmd")
        if err != nil {
            return secret,out, err
        }
        j := JobConfig{cmd: cmd}
        out = append(out,j)
    }
    return secret,out, nil
}
