package main
import (
    "net/http"
    "fmt"
    "log"
    "encoding/json"
)
type Conn struct {
    regex string
}
func (f *Conn) handleWebHook(w http.ResponseWriter, r *http.Request) {
    var gitpayload Gitpayload
    json.NewDecoder(r.Body).Decode(&gitpayload)
    if matchstring(gitpayload.Compare,f.regex) {
        fmt.Println("match!")
    }
    fmt.Println(r.Header["X-Github-Event"][0])
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
