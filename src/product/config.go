package product

import (
    "encoding/json"
    "os"
    "fmt"
)

type Config struct {
    Database    map[string]string
    Mongo       map[string]string
    Redis       map[string]Redis
    Port        int
}

var config Config

func InitConfig(){
    file, err := os.Open("config/conf.json")
    if err != nil {
      fmt.Println("error:", err)
    }
    decoder := json.NewDecoder(file)
    configuration := map[string]Config{}
    err = decoder.Decode(&configuration)
    if err != nil {
      fmt.Println("error:", err)
    }
    fmt.Println(configuration)
    config = configuration["devel"]
}
