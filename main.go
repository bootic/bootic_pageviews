package main

import (
	"bootic_pageviews/request"
	"bootic_pageviews/udp"
	"flag"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

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
		go request.ProcessAndSend(params, query, publisher)
	}

}

func main() {

	var (
		udp_host  string
		http_host string
		gif_path  string
	)

	flag.StringVar(&udp_host, "udphost", "localhost:5555", "UDP host:port to send packets to")
	flag.StringVar(&http_host, "httphost", "localhost:8080", "HTTP host:port to serve tracking gif")
	flag.StringVar(&gif_path, "gifpath", "", "Absolute path to 1x1px tracking gif")

	flag.Parse()

	pub := udp.NewPublisher(udp_host)

	router := mux.NewRouter()

	router.HandleFunc("/r/{app_name}/{account_name}/{type}", PageviewsHandler(gif_path, pub)).Methods("GET")

	http.Handle("/", router)
	log.Println("Starting HTTP endpoint at", http_host)
	log.Fatal(http.ListenAndServe(http_host, nil))
}
