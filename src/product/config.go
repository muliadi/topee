package product

import (
    // "database/sql"
    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
    "gopkg.in/mgo.v2"
    "log"
    "github.com/ruizu/gcfg"
    // "reflect"
    // "time"
    // "strconv"
)

var db *sqlx.DB
var db_product *sqlx.DB
var db_cron *sqlx.DB

var mgo_prod *mgo.Session
var redisconn RedisStruct
var config Config
var BlacklistRule map[string]int

func InitConfig(env string){
  
    filepath := "config/conf-"+env+".ini"
    err := gcfg.ReadFileInto(&config, filepath)
    checkErr(err, "Connect Failed")
    
    InitDb(config.Postgres)
    InitMongo(config.Mongo.Product)
    InitRedis(config.Redis)
    InitBlacklistRule()
}

func InitDb(DBConn PgStruct) {
    dbconn, err := sqlx.Open("postgres", DBConn.Main)
    checkErr(err, "Connect MainDB Failed")
    db = dbconn
    
    dbProd, err := sqlx.Open("postgres", DBConn.Product)
    checkErr(err, "Connect ProductDB Failed")
    db_product = dbProd
    
    dbCron, err := sqlx.Open("postgres", DBConn.Cron)
    checkErr(err, "Connect ProductDB Failed")
    db_cron = dbCron
}

func InitMongo(mongo_product string){
    mongo, err := mgo.Dial(mongo_product)
    checkErr(err, "Fail connect mongo product")
    mongo.SetMode(mgo.Monotonic, true)
    mongo.Ping()
    mgo_prod = mongo
}

func InitRedis(redismap RedisStruct){
    redisconn = redismap
}

func InitBlacklistRule(){
    BlacklistRule = map[string]int{
        "PRD_RULE_BAN_KEYWORD"              : 3,
        "PRD_RULE_WARN_KEYWORD"             : 4,
        "PRD_RULE_BAJAKAN_MAX_HARGA"        : 5,
        "PRD_RULE_BAJAKAN_KATEGORI"         : 6,
        "PRD_RULE_BAJAKAN_BAN_CONDITION"    : 7,
        "PRD_RULE_BAJAKAN_WARN_CONDITION"   : 8,
        "PRD_RULE_CATALOG_BLACKLIST"        : 20,
    }
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}


//==============================================================================
//  STRUCT FOR CONFIG
//==============================================================================
type Config struct {
    Postgres    PgStruct
    Mongo       MongoStruct
    Redis       RedisStruct
}

type PgStruct struct {
    Main        string
    Product     string
    Cron        string
}

type MongoStruct struct{
    Product     string
}

type RedisStruct struct{
    Redis_12_3  string
    Redis_89_5  string
    Redis_89_2  string
    Redis_22_6  string
}
