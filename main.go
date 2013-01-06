package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/url"
    "fmt"
    "bootic_pageviews/udp"
    "os"
)

func PageviewsHandler() (handle func(http.ResponseWriter, *http.Request)) {
  
  // cache the file once
  gif, _ := ioutil.ReadFile("./tiny.gif")
  
  // return closure with access to gif image
  return func(res http.ResponseWriter, req *http.Request) {
    // Write gif response right away
    res.Header().Add("Content-Type", "image/gif")
    res.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, private")
    res.Header().Add("Pragma", "no-cache")
    res.Header().Add("Expires", "Fri, 24 Nov 2000 01:00:00 GMT")

    res.Write(gif)
    
    // Get request data for async processing
    params := mux.Vars(req)
    userAgent := req.UserAgent()
    params["ua"] = userAgent
    query, _ := url.ParseQuery(req.URL.RawQuery)
    // async send data to events collector
    go udp.ProcessAndSend(params, query)
  }
  
}

func main() {
  
  udp_host  := os.Getenv("DATAGRAM_IO_UDP_HOST")
	http_host := os.Getenv("STATS_HTTP_HOST")
	
  udp.Init(udp_host)
  
  router := mux.NewRouter()

  router.HandleFunc("/r/{app_name}/{account_name}/{type}", PageviewsHandler()).Methods("GET")
  
  http.Handle("/", router)
  fmt.Println("Starting HTTP endpoint at", http_host)
  log.Fatal(http.ListenAndServe(http_host, nil))
}