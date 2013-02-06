package main

import (
    "log"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/url"
    "bootic_pageviews/udp"
    "github.com/mssola/user_agent"
    "time"
    "os"
)

func processAndSend(params map[string]string, query url.Values, publisher *udp.Publisher) {
  defer func() {
    if err := recover(); err != nil {
      log.Println("Goroutine failed:", err)
    }
  }()
    
  data := make(map[string]interface{})
  
  data["app"] = params["app_name"]
  data["account"] = params["account_name"]
  
  ua := new(user_agent.UserAgent)
  ua.Parse(params["ua"])
  
  name, version := ua.Browser()
  
  browser := make(map[string]string)
  browser["name"] = name
  browser["version"] = version
  browser["os"] = ua.OS()
  
  data["browser"] = browser
  
  for k, _ := range query {
    if k != "ua" {
      data[k] = query[k][0]
    }
  }
  
  event := make(map[string]interface{})
  event["time"] = time.Now()
  event["type"] = params["type"]
  event["data"] = data

  log.Println("send", event)
  publisher.Publish(event)
}

func PageviewsHandler(gif_path string, publisher *udp.Publisher) (handle func(http.ResponseWriter, *http.Request)) {
  
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
    go processAndSend(params, query, publisher)
  }
  
}

func main() {
  
  udp_host  := os.Getenv("DATAGRAM_IO_UDP_HOST")
	http_host := os.Getenv("STATS_HTTP_HOST")
	gif_path  := os.Getenv("GIF_PATH")
	
  pub := udp.NewPublisher(udp_host)
  
  router := mux.NewRouter()

  router.HandleFunc("/r/{app_name}/{account_name}/{type}", PageviewsHandler(gif_path, pub)).Methods("GET")
  
  http.Handle("/", router)
  log.Println("Starting HTTP endpoint at", http_host)
  log.Fatal(http.ListenAndServe(http_host, nil))
}