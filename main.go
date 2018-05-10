package main
import (
    "net/http"
    "fmt"
    "log"
    "encoding/json"
)
func handleWebHook(w http.ResponseWriter, r *http.Request) {
    config,_ := readconfig()
    var gitpayload Gitpayload
    json.NewDecoder(r.Body).Decode(&gitpayload)
    if matchstring(gitpayload.Compare,config[2]) {
        fmt.Println("match!")
    }
    fmt.Println(r.Header["X-Github-Event"][0])
}
func main() {
    config,_ := readconfig()
    http.HandleFunc("/", handleWebHook)
    err := http.ListenAndServeTLS(config[4]+":"+config[5],config[0],config[1], nil)
    if err != nil {
        log.Fatal("Unable to Listen to port", err)
    }
}
