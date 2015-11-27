package product

import(
    "bytes"
    // "fmt"
    // "time"
    "strconv"
    // "regexp"
    
    // "github.com/extemporalgenome/slug"
    "gopkg.in/redis.v3"
    // "gopkg.in/mgo.v2/bson"
)

func UpdateProduct(input *ProductInput, current *ProductInput){
    
    //set initial status
    input.Status = current.Status
    
    //filter blacklist if desc is updated
    if input.ShortDesc != ""{
        if res, word:=CheckBlacklist(input.ShortDesc, BlacklistRule["PRD_RULE_BAN_KEYWORD"]); res == true {
            input.Status = -2
            input.ProductStatus = "Banned"
            input.PendingReason = "Warned Product Description because keyword " + word
            input.PendingStatus = 2
        } else if res, word:=CheckBlacklist(input.ShortDesc, BlacklistRule["PRD_RULE_WARN_KEYWORD"]); res == true {
            input.Status = -1
            input.ProductStatus = "Warned"
            input.PendingReason = "Warned Product Description because keyword " + word
            input.PendingStatus = 2
        }
    }

    //check pirated product
    //if category or price is updated, check for piracy
    if input.ChildCatId != 0 && input.NormalPrice != 0{ 
        //if both changed
        //pirated product filter
        if CheckBajakan(input.NormalPrice, input.ChildCatId) == true {
            input.Status = -1
            input.ProductStatus = "Warned"
            input.PendingReason = "Warned because suspected as pirated product"
            input.PendingStatus = 3
        }
    } else if input.ChildCatId == 0 && input.NormalPrice != 0{
        //if only pridce
        if CheckBajakan(input.NormalPrice, current.ChildCatId) == true {
            input.Status = -1
            input.ProductStatus = "Warned"
            input.PendingReason = "Warned because suspected as pirated product"
            input.PendingStatus = 3
        }
    } else if input.ChildCatId != 0 && input.NormalPrice == 0 {
        //if only category changed
        if CheckBajakan(current.NormalPrice, input.ChildCatId) == true {
            input.Status = -1
            input.ProductStatus = "Warned"
            input.PendingReason = "Warned because suspected as pirated product"
            input.PendingStatus = 3
        }
    }
        
    //set product status
    SetProductStatus(input)
    
    //upsert product pending reason
    if input.Status == -1 || input.Status == -2 {
        UpsertPendingReason(input)
    }
    
    if catalog is changed
    if input.AddToCatalog == 0{
        //remove from catalog
        RemoveFromCatalog(input.ProductId)
    } else {
        //move to another catalog
        RemoveFromCatalog(input.ProductId)
        res_ctg_prd_desc, _ := CheckBlacklist(input.ShortDesc, BlacklistRule["PRD_RULE_CATALOG_BLACKLIST"])
        if res_ctg_prd_desc==false {
            AddToCatalog(input)
            UpdateCron(input.ProductId, "price_alert_catalog")
        }
    }
    
    //add to etalase
    if input.AddToEtalase == 1 && input.Status != -1 && input.Status != -2 && input.Status != 3 {
        AddToEtalase(input.ProductId, input.EtalaseId)
        UpdateCron(input.ProductId, "price_alert_product")
    }
    
    //update wholesale price
    if len(input.Wholesale) > 0{
        AddWholesalePrice(input)
    }
    
    //update product returnable
    if input.Returnable != -1 {
        UpsertReturnable(input.ProductId, input.ShopId, input.Returnable)
    }
    
    //insert product data to mongoDB
    UpsertProductList(input)
    
    //add product history
    AddProductHistory(input)
    
    //scan any phone number in description and insert to mongo for security team to use it later
    ScanPhoneNumber(input.ShortDesc, input.ProductId)
    
    UpdateWsProduct(input, current)
    
    delete_redis("dir_product:view_list:p_"+strconv.FormatInt(input.ProductId, 16))
    delete_redis("dir_product:view_gallery:p_"+strconv.FormatInt(input.ProductId, 16))
    delete_redis("class:product:p_"+strconv.FormatInt(input.ProductId, 16))

}


func RemoveFromCatalog(prod_id int64){
    buff := bytes.NewBufferString(`
        UPDATE ws_catalog_product 
        SET 
            status = 0,
            update_time = $1
        WHERE
            product_id = $2
            AND status = 1
    `)
    
    query := db.Rebind(buff.String())
    db.MustExec(
        query,
        Now(),
        prod_id)
}

func UpdateCron(prod_id int64, cron_type string){
    buff := bytes.NewBufferString(`
        UPDATE ws_cron_job
        SET
            status = 1,
            create_time = $1,
            update_time = null
        WHERE 
            id = $2
            AND type = $3
    `)
    
    query := db_cron.Rebind(buff.String())
    db_cron.MustExec(query, Now(), prod_id, cron_type)
}


func UpdateWsProduct(input *ProductInput, current *ProductInput){
    
    //if category not updated 
    if input.ChildCatId == 0{
        input.ChildCatId = current.ChildCatId
    }
    
    //if min order not updated
    if input.MinOrder == 0{
        input.MinOrder = current.MinOrder
    }
    
    //if price currency not updated
    if input.PriceCurrency == 0{
        input.PriceCurrency = current.PriceCurrency
    }
    
    //if price not updated
    if input.NormalPrice == 0{
        input.NormalPrice = current.NormalPrice
    } 
    
    //if any price updated, set last updated price
    if input.PriceCurrency == 0 || input.NormalPrice == 0{
        input.LastPrcUpdate = current.LastPrcUpdate
    } else {
        input.LastPrcUpdate = Now()
    }
    
    if input.Weight.Unit == 0 {
        input.Weight.Unit = current.Weight.Unit
    }
    
    if input.Weight.Numeric == 0{
        input.Weight.Numeric = current.Weight.Numeric
    }
    
    if input.Insurance == -1{
        input.Insurance = current.Insurance
    }
    
    if input.ShortDesc == ""{
        input.ShortDesc = current.ShortDesc
    }
    
    if input.Condition == 0{
        input.Condition = current.Condition
    }
    
    if input.Returnable == -1{
        input.Returnable = current.Returnable
    }
    
    buff := bytes.NewBufferString(`
        UPDATE ws_product 
        SET 
            status          = $1,
            child_cat_id    = $2,
            min_order       = $3,
            price_currency  = $4,
            normal_price    = $5,
            weight_unit     = $6,
            weight          = $7,
            must_insurance  = $8,
            short_desc      = $9,
            condition       = $10,
            update_solr     = $11
        WHERE
            product_id = $12
    `)
    
    query := db.Rebind(buff.String())
    _, err := db.Query(
        query,
        input.Status,
        input.ChildCatId,
        input.MinOrder,
        input.PriceCurrency,
        input.NormalPrice,
        input.Weight.Unit,
        input.Weight.Numeric,
        input.Insurance,
        input.ShortDesc,
        input.Condition,
        Now(),
        input.ProductId)
    
    checkErr(err, "fail update")
}


func delete_redis(key string){
     rds := redis.NewClient(&redis.Options{
        Addr        : redisconn.Redis_89_5,
        Password    : "", // no password set
        DB          : 0,  // use default DB
    })
   
    rds.Del(key)
}
