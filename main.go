package main

import (
    "log"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/url"
    "bootic_pageviews/udp"
    "os"
)

func PageviewsHandler(gif_path string) (handle func(http.ResponseWriter, *http.Request)) {
  
  // cache the file once
  gif, _ := ioutil.ReadFile(gif_path)
  
  // return closure with access to gif image
  return func(res http.ResponseWriter, req *http.Request) {
    res.Header().Add("Content-Type", "image/gif")
    res.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, private, proxy-revalidate")
    res.Header().Add("Pragma", "no-cache")
    res.Header().Add("Expires", "Fri, 24 Nov 2000 01:00:00 GMT")

    // Content-Length. This is not a streaming connection.
    res.Header().Add("Content-Length", strconv.Itoa(len(gif)))
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
	gif_path  := os.Getenv("GIF_PATH")
	
  udp.Init(udp_host)
  
  router := mux.NewRouter()

  router.HandleFunc("/r/{app_name}/{account_name}/{type}", PageviewsHandler(gif_path)).Methods("GET")
  
  http.Handle("/", router)
  log.Println("Starting HTTP endpoint at", http_host)
  log.Fatal(http.ListenAndServe(http_host, nil))
}