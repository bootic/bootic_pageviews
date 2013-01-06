package udp

import (
  "net"
  "net/url"
  "fmt"
  "time"
  "encoding/json"
)

type Event struct {
	Type string   `json:"type"`
	Time time.Time `json:"time"`
	Data map[string]string `json:"data"`
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
  
  data := make(map[string]string)
  
  data["app"] = params["app_name"]
  data["account"] = params["account_name"]
  
  for k, _ := range query {
    data[k] = query[k][0]
  }
  
  event := Event{
    Time: time.Now(),
    Type: params["type"],
    Data: data,
  }
  
  json, err := json.Marshal(event)
  
  if err != nil {
    fmt.Println("Could not marshal JSON:", err)
  }
  
  udpConn.Write(json)
  
  // fmt.Println("sending data" + params["app_name"] + " lalalal " + params["ua"])
}