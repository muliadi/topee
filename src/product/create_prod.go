package product

import(
    "bytes"
    "fmt"
    "time"
    "strconv"
    
    "github.com/extemporalgenome/slug"
    "gopkg.in/redis.v3"
    "gopkg.in/mgo.v2/bson"
)

func CreateProduct(input *ProductInput){
        
    //filter blacklist
    if res, word:=CheckBlacklist(input.ProductName, BlacklistRule["PRD_RULE_BAN_KEYWORD"]); res == true {
        input.Status = -2
        input.PendingReason = "Warned Product Description because keyword " + word
        input.PendingStatus = 1
    } else if res, word:=CheckBlacklist(input.ShortDesc, BlacklistRule["PRD_RULE_BAN_KEYWORD"]); res == true {
        input.Status = -2
        input.PendingReason = "Warned Product Description because keyword " + word
        input.PendingStatus = 2
    } else if res, word:=CheckBlacklist(input.ProductName, BlacklistRule["PRD_RULE_WARN_KEYWORD"]); res == true {
        input.Status = -1
        input.PendingReason = "Warned Product Description because keyword " + word
        input.PendingStatus = 1
    } else if res, word:=CheckBlacklist(input.ShortDesc, BlacklistRule["PRD_RULE_WARN_KEYWORD"]); res == true {
        input.Status = -1
        input.PendingReason = "Warned Product Description because keyword " + word
        input.PendingStatus = 2
    }
    
    //insert product into ws_product
    input.ProductId = InsertProduct(input)
    
    //add to etalase
    if input.AddToEtalase == 1 && input.Status != -1 && input.Status != -2 && input.Status != 3 {
        AddToEtalase(input.ProductId, input.EtalaseId)
    }
    
    //add to catalog
    //first check blacklist for catalog
    res_ctg_prd_desc, _ := CheckBlacklist(input.ShortDesc, BlacklistRule["PRD_RULE_CATALOG_BLACKLIST"])
    res_ctg_prd_name, _ := CheckBlacklist(input.ShortDesc, BlacklistRule["PRD_RULE_CATALOG_BLACKLIST"])
    if res_ctg_prd_name==false && res_ctg_prd_desc==false {
        AddToCatalog(input)
    }
    
    //create product alias
    CreateAlias(input.ProductId, input.ProductName, input.ShopId)
    
    //insert wholesale price
    if len(input.Wholesale) > 0{
        AddWholesalePrice(input)
    }
    
    //update current shop max position
    UpdateMaxPosition(input.ShopId, input.Position)
    
    //set product returnable
    UpsertReturnable(input.ProductId, input.ShopId, input.Returnable)
    
    //add to redis sitemap
    AddSitemapProduct(input.ProductId)
    
    //upsert product pending reason
    if input.Status == -1 || input.Status == -2 {
        UpsertPendingReason(input)
    }
    
    //set product stat redis
    var datamap = map[string]string{
        "count_review"      : "0",
        "count_talk"        : "0",
        "count_sold"        : "0",
        "rating"            : "0",
        "count_tx_success"  : "0",
        "count_tx_reject"   : "0",
        "count_view"        : "0",
    }
    SetProductStatRedis(input.ProductId, datamap)
    
    //insert product data to mongoDB
    UpsertProductList(input)
    
    AddBroadcast(input.ProductId, input.ShopId)
}

// MAX POSITION FUNCTION - START
func GetMaxPosition(shop_id int64) int64{
    
    fmt.Println("")
    
    var max_position int64
    row := db.QueryRow("SELECT max_position FROM ws_shop_max_position WHERE shop_id=$1", shop_id)
    row.Scan(&max_position)
    return max_position
}

func UpdateMaxPosition(shop_id int64, pos int64){
    current_pos := GetMaxPosition(shop_id)
    if current_pos == 0{
        db.Exec(`INSERT INTO 
                    ws_shop_max_position 
                (shop_id, max_position)
                VALUES ($1, $2)`, shop_id, pos)
    } else {
        db.Exec(`UPDATE ws_shop_max_position
                SET max_position = $1
                WHERE shop_id = $2`, pos, shop_id)
    }
}
// MAX POSITION FUNCTION - END

func InsertProduct(product *ProductInput) int64{
    //get max position in current shop
    product.Position = GetMaxPosition(product.ShopId)
    product.Position++
    
    var product_id int64
    now := Now()
    
    err := db.QueryRow(`
        INSERT INTO ws_product
        (
            shop_id,
            child_cat_id,
            product_name,
            short_desc,
            normal_price,
            price_currency,
            status,
            weight,
            weight_unit,
            min_order,
            create_by,
            create_time,
            position,
            must_insurance,
            last_update_price,
            condition,
            update_solr
        )
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
        RETURNING product_id`,
        product.ShopId,
        product.ChildCatId,
        product.ProductName,
        product.ShortDesc,
        product.NormalPrice,
        product.PriceCurrency,
        product.Status,
        product.Weight.Numeric,
        product.Weight.Unit,
        product.MinOrder,
        product.UserId,
        now,
        product.Position,
        product.Insurance,
        now,
        product.Condition,
        now).Scan(&product_id)

    if err != nil {
        checkErr(err, "fail Insert")
        return 0
    } else {
        return product_id
    }
}

func AddToEtalase(prod_id int64, menu_id int64){
    buff := bytes.NewBufferString(`
        INSERT INTO ws_product_menu (
            product_id,
            menu_id
        ) VALUES ($1, $2)
    `)
    
    query := db.Rebind(buff.String())
    db.MustExec(query, prod_id, menu_id)
}

func AddToCatalog(product *ProductInput){
    buff := bytes.NewBufferString(`
        INSERT INTO ws_catalog_product (
            ctg_id,
            product_id,
            shop_id,
            status,
            create_by,
            create_time
        ) VALUES ($1, $2, $3, $4, $5, $6)
    `)
    
    query := db.Rebind(buff.String())
    db.MustExec(
        query,
        product.CatalogId,
        product.ProductId,
        product.ShopId,
        1,
        product.UserId,
        Now())
}

func CreateAlias(prod_id int64, prod_name string, shop_id int64){
    //check alias
    var count int
    var loop int
    loop = 0
    
    var found int
    found = 0
    
    var key string
    plain_key := slug.Slug(prod_name)
    
    for found == 0{
        
        if loop == 0{
            key = plain_key
        } else {
            key = plain_key + "-" + strconv.Itoa(loop)
        }
            
        row := db_product.QueryRow(`
            SELECT 
                count(product_id)
            FROM 
                ws_product_alias 
            WHERE 
                product_key = $1
                AND shop_id = $2`,
            key,
            shop_id)
        row.Scan(&count)
        
        if(count == 0){
            found = 1
        } else {
            loop++
        }
    }
    
    //insert alias to postgre
    _, err := db_product.Exec(`
        INSERT INTO ws_product_alias
            (product_id, product_key, shop_id)
        VALUES($1, $2, $3)
    `, prod_id, key, shop_id)    
    checkErr(err, "Fail create alias in postgres")
    
    //insert alias to mongodb
    cmgo := mgo_prod.DB("product_dev").C("product_alias")
    alias := &ProductAlias{
        ProductId   : prod_id,
        ProductKey  : key,
        ShopId      : shop_id,
    }
    _, err = cmgo.Upsert(bson.M{"product_id":prod_id, "shop_id":shop_id}, alias)
    checkErr(err, "Fail create alias in mongodb")
    
    //delete product alias in redis
    rds := redis.NewClient(&redis.Options{
        Addr        : redisconn.Redis_12_3,
        Password    : "", // no password set
        DB          : 0,  // use default DB
    })
    
    rds.Del("svq:aliasing-id_product-"+key+"-"+string(shop_id))
    rds.Close()
}

func AddWholesalePrice(product *ProductInput){
    wholesale := make(map[string]int64)
    
    var prd_prc_id int64
    var loop int = 1
    for _, ws := range product.Wholesale {
        wholesale["qty_min_"+strconv.Itoa(loop)] = ws.Min
        wholesale["qty_max_"+strconv.Itoa(loop)] = ws.Max
        wholesale["prd_prc_"+strconv.Itoa(loop)] = ws.Price
        loop++
    }
    
    now := Now()
    
    //insert into ws_prd_prc
    err := db.QueryRow(`
        INSERT INTO ws_prd_prc
            (
                product_id,
                create_by,
                create_time,
                qty_min_1,
                qty_max_1,
                prd_prc_1,
                qty_min_2,
                qty_max_2,
                prd_prc_2,
                qty_min_3,
                qty_max_3,
                prd_prc_3,
                qty_min_4,
                qty_max_4,
                prd_prc_4,
                qty_min_5,
                qty_max_5,
                prd_prc_5
            ) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
        RETURNING prd_prc_id`,
        product.ProductId,
        product.UserId,
        now,
        wholesale["qty_min_1"],
        wholesale["qty_max_1"],
        wholesale["prd_prc_1"],
        wholesale["qty_min_2"],
        wholesale["qty_max_2"],
        wholesale["prd_prc_2"],
        wholesale["qty_min_3"],
        wholesale["qty_max_3"],
        wholesale["prd_prc_3"],
        wholesale["qty_min_4"],
        wholesale["qty_max_4"],
        wholesale["prd_prc_4"],
        wholesale["qty_min_5"],
        wholesale["qty_max_5"],
        wholesale["prd_prc_5"]).Scan(&prd_prc_id)
    
    if err != nil {
        checkErr(err, "fail Insert wholesale in postgres")
    } else {
        UpsertWholesaleMongo(product.UserId, product.ProductId, prd_prc_id, wholesale)
    }
}

func UpsertWholesaleMongo(user_id int64, prod_id int64, prd_prc_id int64, wholesale map[string]int64){
    cmgo := mgo_prod.DB("product_dev").C("product_price_history")
    wsMongo := WholesaleMongo{
        UpdateTime  : time.Now().Unix(),
        UpdateBy    : user_id,
        ProductId   : prod_id,
    }
    
    if wholesale["qty_min_1"] > 0{
        qty_min_1 := wholesale["qty_min_1"]
        qty_max_1 := wholesale["qty_max_1"]
        prd_prc_1 := wholesale["prd_prc_1"]
        wsMongo.QtyMin1 = qty_min_1
        wsMongo.QtyMax1 = qty_max_1
        wsMongo.PrdPrc1 = prd_prc_1
    }
    
    if wholesale["qty_min_2"] > 0{
        qty_min_2 := wholesale["qty_min_2"]
        qty_max_2 := wholesale["qty_max_2"]
        prd_prc_2 := wholesale["prd_prc_2"]
        wsMongo.QtyMin2 = qty_min_2
        wsMongo.QtyMax2 = qty_max_2
        wsMongo.PrdPrc2 = prd_prc_2
    }
    
    if wholesale["qty_min_3"] > 0{
        qty_min_3 := wholesale["qty_min_3"]
        qty_max_3 := wholesale["qty_max_3"]
        prd_prc_3 := wholesale["prd_prc_3"]
        wsMongo.QtyMin3 = qty_min_3
        wsMongo.QtyMax3 = qty_max_3
        wsMongo.PrdPrc3 = prd_prc_3
    }
    
    if wholesale["qty_min_4"] > 0{
        qty_min_4 := wholesale["qty_min_4"]
        qty_max_4 := wholesale["qty_max_4"]
        prd_prc_4 := wholesale["prd_prc_4"]
        wsMongo.QtyMin4 = qty_min_4
        wsMongo.QtyMax4 = qty_max_4
        wsMongo.PrdPrc4 = prd_prc_4
    }
    
    if wholesale["qty_min_5"] > 0{
        qty_min_5 := wholesale["qty_min_5"]
        qty_max_5 := wholesale["qty_max_5"]
        prd_prc_5 := wholesale["prd_prc_5"]
        wsMongo.QtyMin5 = qty_min_5
        wsMongo.QtyMax5 = qty_max_5
        wsMongo.PrdPrc5 = prd_prc_5
    }
    
    _, err := cmgo.Upsert(bson.M{"product_id":prod_id}, wsMongo)
    checkErr(err, "Fail create wholesale in mongodb")
    
    //delete product wholesale price in redis
    rds := redis.NewClient(&redis.Options{
        Addr        : redisconn.Redis_12_3,
        Password    : "", // no password set
        DB          : 0,  // use default DB
    })
    
    rds.Del("lib_cache:wholesale_price")
    rds.Del("lib_cache:last_update_price")
    rds.Del("lib_cache:wholesale_update_time")
    rds.Del("lib_cache:facade:product:get_wholesale_price:"+string(prod_id))
    rds.Close()
}

func UpsertReturnable(prod_id int64, shop_id int64, returnable int64){
    cmgo := mgo_prod.DB("product_dev").C("product_info")
    prdinfo := &ProductInfoMongo{
        Returnable  : returnable,
        ShopId      : shop_id,
        ProductId   : prod_id,
    }
    
    _, err := cmgo.Upsert(bson.M{"product_id": prod_id, "shop_id":shop_id}, prdinfo)
    checkErr(err, "Fail create product_info in mongodb")
}

func UpsertProductList(product *ProductInput){
    cmgo := mgo_prod.DB("tokopedia_product_dev").C("product_list")
    prdlist := &ProductListMongo{
        ChildCatId      : product.ChildCatId,
        Condition       : product.Condition,
        CreateTime      : time.Now().Unix(),
        CatalogId       : product.CatalogId,
        MenuId          : product.EtalaseId,
        NormalPrice     : product.NormalPrice,
        Position        : product.Position,
        PriceCurrency   : product.PriceCurrency,
        ProductId       : product.ProductId,
        ProductName     : product.ProductName,
        RupiahPrice     : product.NormalPrice,
        ShopId          : product.ShopId,
        UpdateTime      : time.Now().Unix(),
        Status          : product.Status,
    }
    
    _, err := cmgo.Upsert(bson.M{"product_id": product.ProductId, "shop_id":product.ShopId}, prdlist)
    checkErr(err, "Fail to upsert product_list in mongodb")
}

func AddSitemapProduct(prod_id int64){
    rds := redis.NewClient(&redis.Options{
        Addr        : redisconn.Redis_89_5,
        Password    : "", // no password set
        DB          : 0,  // use default DB
    })
    
    value := redis.Z{
        Score       : float64(prod_id),
        Member      : prod_id,
    }
    
    rds.ZAdd("sitemap:product", value)
}

func InsertCron(prod_id int64, cron_type string){
    buff := bytes.NewBufferString(`
        INSERT INTO ws_cron_job (
            id,
            type,
            status,
            create_time
        ) VALUES ($1, $2, 1, $3)
    `)
    
    query := db_cron.Rebind(buff.String())
    db_cron.MustExec(query, prod_id, cron_type, Now())
}

func SetProductStatRedis(prod_id int64, data map[string]string){
    rds := redis.NewClient(&redis.Options{
        Addr        : redisconn.Redis_89_2,
        Password    : "", // no password set
        DB          : 0,  // use default DB
    })
    
    var key string = "product_stats_hash:"+strconv.FormatInt(prod_id, 16)
    for field, value := range data {
        rds.HMSet(key, field, value)
    }
}

func AddBroadcast(prod_id int64, shop_id int64){
    rds := redis.NewClient(&redis.Options{
        Addr        : redisconn.Redis_22_6,
        Password    : "", // no password set
        DB          : 0,  // use default DB
    })
    
    //list:feed_product
    value := redis.Z{
        Score       : float64(time.Now().Unix()),
        Member      : prod_id,
    }
    rds.ZAdd("list:feed_product:"+strconv.FormatInt(shop_id, 16), value)
    
    //fave_product_broadcast
    keymap := map[int]string{
        1 : "fave_product_broadcast:p1",
        2 : "fave_product_broadcast:p2",
        3 : "fave_product_broadcast:p3",
    }
    key := keymap[(int(prod_id) % 3)+1]
    
    shop_prod := strconv.FormatInt(shop_id, 16)+"-"+strconv.FormatInt(prod_id, 16)
    rds.LRem(key, 0, shop_prod)
    rds.LPush(key, shop_prod)
}

func Now() string{
    // time.LoadLocation("Asia/Jakarta")
    // t := time.Now()
    // now := t.Format("20060102150405")
    // year := now[0:4]
    // month := now[4:6]
    // day := now[6:8]
    // hour := now[8:10]
    // min := now[10:12]
    // sec := now[12:14]
    
    // var timestamp string
    // timestamp = year + "-" + month + "-" + day + " " + hour + ":" + min + ":" + sec
    // return timestamp
    now := time.Now()
    secs := now.Unix()
    timestamp := time.Unix(secs, 0).String()
    result := timestamp[0:19]
    return result
}
