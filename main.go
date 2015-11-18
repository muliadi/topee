package main

import (
    "github.com/gin-gonic/gin"
    "github.com/alvinantonius/produp/src/product"
    // "database/sql"
    // _ "github.com/lib/pq"
    // "github.com/jmoiron/sqlx"
    "log"
    // "fmt"
    // "time"
    // "strconv"
)

func main(){
    app := gin.Default()
    app.GET("/", product.Index)
    app.POST("/create", product.Create)
    app.Run(":8000")
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
