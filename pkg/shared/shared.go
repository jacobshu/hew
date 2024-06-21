package shared

import (
  "encoding/json"
  "fmt"
)

func Pprint(data interface{}) {
  s, err := json.MarshalIndent(data, "", "\t")
  if err != nil {
    fmt.Printf("Pprint error\n%#v\n", err)
  }
  fmt.Println(string(s))
}
