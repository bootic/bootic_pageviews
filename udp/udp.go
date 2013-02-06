package udp

import (
  "net"
  "log"
  "encoding/json"
)

var udpConn *net.UDPConn

type Publisher struct {
  conn *net.UDPConn
}

func NewPublisher(hostAndPort string) (*Publisher) {
  udpAddr, err := net.ResolveUDPAddr("udp", hostAndPort)
  if err != nil { panic("Could not connect to UDP server") }
  var err2 error
  udpConn, err2 = net.DialUDP("udp", nil, udpAddr) 
  if err2 != nil { panic(err2) }
  
  pub := &Publisher{
    conn: udpConn,
  }
  
  return pub
}

func (p *Publisher) Publish(event map[string]interface{}) {
  json, err := json.Marshal(event)
  
  if err != nil {
    log.Println("Could not marshal JSON:", err)
  }
  
  p.conn.Write(json)
}