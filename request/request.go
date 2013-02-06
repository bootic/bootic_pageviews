package request

import(
  "net/url"
  "bootic_pageviews/udp"
  "log"
  "github.com/mssola/user_agent"
  "time"
)

func ProcessAndSend(params map[string]string, query url.Values, publisher *udp.Publisher) {
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

  publisher.Publish(event)
}