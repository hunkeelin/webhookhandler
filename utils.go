package main

import (
    "strings"
    "github.com/theckman/go-flock"
    "log"
    "sync"
    "os"
    "os/exec"
    "regexp"
    "fmt"
    "time"
)

func waitforqueue(dir string) *flock.Flock{
    fileLock := flock.NewFlock(dir)

    locked, err := fileLock.TryLock()

    if err != nil {
        log.Fatal("unable to lock the file at ",dir)
    }

    if locked {
        return fileLock
    } else {
        fmt.Println("wait 10 seconds there's another process running")
        time.Sleep(10 * time.Second)
        return waitforqueue(dir)
    }
}

func matchstring(s,regex string) bool {
    match, err := regexp.MatchString(regex, s)
    fmt.Println("")
    if err != nil {
        log.Fatal("regex matching problem")
    }
    return match
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func exist(path string)(bool, error){
    _,err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return true, err
}

func runshell(cmd string, args []string) error{
    err := exec.Command(cmd, args...).Run()
  //  acmd := exec.Command(cmd, args...)
  //  out,err := acmd.CombinedOutput()
  //  fmt.Printf("%s\n",out)
    return err
}

func removestring(s []string, pattern string) []string{
    var toreturn []string
    for _,raw_element := range s {
        element := strings.Replace(raw_element," ","",-1)
        if strings.HasPrefix(element,pattern){
            ele := strings.Replace(element,pattern,"",-1)
            if strings.HasPrefix(ele,"HEAD") || strings.HasPrefix(ele,"master"){
                continue
            }
            toreturn = append(toreturn,ele)
        }
    }
    return toreturn
}

func outshell(cmd string, args []string) (string,error){
    output, err := exec.Command(cmd, args...).Output()
    return string(output), err
}
func cleandir(s,env []string, workers int){
    sema := make(chan struct{}, workers)
    wg := sync.WaitGroup{}
    for _,element := range s {
        if !stringInSlice(element,env){
            sema <- struct{}{}
            wg.Add(1)
            go func(element string){
                os.RemoveAll(element)
                <-sema
                wg.Done()
            }(element)
        }
    }
    wg.Wait()
}

func createdir(s,env []string,workers int){
    sema := make(chan struct{}, workers)
    wg := sync.WaitGroup{}
    for _,element := range env {
        if stringInSlice(element,s) == false {
            sema <- struct{}{}
            wg.Add(1)
            go func (element string){
                os.MkdirAll(element+"/"+"modules",0755)
                <-sema
                wg.Done()
            }(element)
        }
    }
    wg.Wait()
}
func checkstring(s,pattern string) {
    if strings.HasPrefix(s, "mod") {
        if string(s[len(s)-1]) != "," {
            log.Fatal("error missing comma on line: ",s)
        }
    }

    checknum := 0

    for _,r := range s {
        c := string(r)
        if c == "'" {
            checknum = checknum + 1
        }
    }

    if checknum !=2 {
        log.Fatal("error missing single quotes on line: ",s)
    }
}

func trim(x string) string {
    pattern := "'"
    checkstring(x,pattern)
    bra := strings.Index(x, pattern)
    if bra < 0 {
        return ""
    }
    rx := x[bra+1:]
    ket := strings.Index(rx, pattern)
    return rx[:ket]
}