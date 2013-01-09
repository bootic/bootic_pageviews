package udp

import (
  "net"
  "net/url"
  "log"
  "time"
  "encoding/json"
  "github.com/mssola/user_agent"
)

type Event struct {
	Type string   `json:"type"`
	Time time.Time `json:"time"`
	Data map[string]interface{} `json:"data"`
}

var udpConn *net.UDPConn

func Init(hostAndPort string) {
  udpAddr, err := net.ResolveUDPAddr("udp", hostAndPort)
  if err != nil { panic("Could not connect to UDP server") }
  var err2 error
  udpConn, err2 = net.DialUDP("udp", nil, udpAddr) 
  if err2 != nil { panic(err2) }
}

func ProcessAndSend(params map[string]string, query url.Values) {
  
  defer func() {
    if err := recover(); err != nil {
      log.Println("Goroutine failed:", err)
    }
  }()
    
  data := make(map[string]interface{})
  
  data["app"] = params["app_name"]
  data["account"] = params["account_name"]
  
  ua := new(user_agent.UserAgent)
  ua.Parse(query["ua"][0])
  
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
  
  event := Event{
    Time: time.Now(),
    Type: params["type"],
    Data: data,
  }
  
  json, err := json.Marshal(event)
  
  if err != nil {
    log.Println("Could not marshal JSON:", err)
  }
  
  udpConn.Write(json)
}