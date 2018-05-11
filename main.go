package main
import (
    "net/http"
    "fmt"
    "log"
    "encoding/json"
)
func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
    var gitpayload Gitpayload
    err := json.NewDecoder(r.Body).Decode(&gitpayload)
    if err != nil {
        log.Fatal("cannot decode body with gitstruct")
    }
   // if matchstring(gitpayload.Compare,f.regex) {
   //     fmt.Println("match!")
   // }
    if gitpayload.doesmatchbody(f.regex) {
        fmt.Println("match")
    } else {
        fmt.Println("not match")
    }
    fmt.Println(r.Header["X-Github-Event"][0])
    fmt.Println(r.Header)
}
func main() {
    newcon := new(Conn)
    config,_ := readconfig()
    newcon.regex = config[2]
    http.HandleFunc("/", newcon.handleWebHook)
    err := http.ListenAndServeTLS(config[4]+":"+config[5],config[0],config[1], nil)
    if err != nil {
        log.Fatal("Unable to Listen to port", err)
    }
}
