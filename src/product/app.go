package product

import (
    // "database/sql"
    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
    "gopkg.in/mgo.v2"
    "log"
    "fmt"
    // "reflect"
    // "time"
    // "strconv"
)

var db *sqlx.DB
var db_product *sqlx.DB
var mgo_prod *mgo.Session
var redisconn *map[string]Redis

func init(){
    InitConfig()
    InitDb(config.Database["main"], config.Database["product"])
    InitMongo(config.Mongo["product"])
    InitRedis(config.Redis)
}

func InitDb(mainDB string, productDB string) {
    dbconn, err := sqlx.Open("postgres", mainDB)
    checkErr(err, "Connect Failed")
    db = dbconn
    
    dbProd, err := sqlx.Open("postgres", productDB)
    checkErr(err, "Connect Failed")
    db_product = dbProd
}

func InitMongo(mongo_product string){
    mongo, err := mgo.Dial(mongo_product)
    checkErr(err, "Fail connect mongo product")
    mongo.Ping()
    mgo_prod = mongo
    fmt.Println(mongo_product)
}

func InitRedis(redismap map[string]Redis){
    redisconn = &redismap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
