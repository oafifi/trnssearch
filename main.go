package main

import (
  "github.com/oafifi/trnssearch/api"
  "fmt"
)

func main()  {

  fmt.Println("Starting server ...")
  api.HandleRequests()

}
