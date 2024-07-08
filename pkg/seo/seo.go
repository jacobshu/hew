package seo

import (
   "bytes"
   "encoding/json"
   "io"
   "log"
   "net/http"
)

func main() {
  postArray := []map[string]string
  postArray = append(postArray, make(map[string]string{
    "target": "dataforseo.com",
    "max_crawl_pages": 10,
    "load_resources": true,
    "enable_javascript": true,
    "custom_js": "meta = {}; meta.url = document.URL; meta;",
    "tag": "some_string_123",
    "pingback_url": "https://your-server.com/pingscript?id=$id&tag=$tag",
  }))
  
  {
    method: 'post',
    url: 'https://api.dataforseo.com/v3/on_page/task_post',
    auth: {
      username: 'login',
      password: 'password'
    },
    data: post_array,
    headers: {
      'content-type': 'application/json'
    }
  }

   postBody, _ := json.Marshal(map[string]string{
      "name":  "Toby",
      "email": "Toby@example.com",
   })
   responseBody := bytes.NewBuffer(postBody)
  
   resp, err := http.Post("https://postman-echo.com/post", "application/json", responseBody)
   if err != nil {
      log.Fatalf("An Error Occured %v", err)
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
   sb := string(body)
   log.Printf(sb)
}
