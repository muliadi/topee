package main

import (
    "github.com/ruizu/gcfg"
)

type Config struct {
    Server        ServerConf
}

type ServerConf struct{
    Host        string
    Env         string
}

func ReadConfig(c *Config, filepath string) bool {
    if err := gcfg.ReadFileInto(c, filepath); err != nil {
        return false
    }
    return true
}
