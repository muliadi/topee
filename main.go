package main

import (
    "github.com/alvinantonius/produp/src/product"
    
    "github.com/julienschmidt/httprouter"
    "net/http"
    au "github.com/ruizu/api-utils"

    "log"
    "log/syslog"
    // "fmt"
    // "time"
    // "strconv"
)

var config Config

func init(){
    ok := ReadConfig(&config, "config/server.ini")
    if !ok {
        log.Fatal("Could not find configuration file")
    }
    
    product.InitConfig(config.Server.Env)
    
    if config.Server.Env == "production" {
        logger, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_DAEMON, "product")
        if err != nil {
            log.Fatal(err)
        }
        log.SetOutput(logger)
        log.SetFlags(0)
    }
}

func main(){
    router := httprouter.New()
    router.NotFound = au.NotFoundHandler
    router.MethodNotAllowed = au.MethodNotAllowedHandler
    router.PanicHandler = au.PanicHandler
    au.Debug = true;
    
    router.GET("/ping", product.Ping)
    
    router.GET("/", product.IndexHandler)
    router.POST("/create", product.CreateHandler)
    
    log.Fatal(http.ListenAndServe(config.Server.Host, router))
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
