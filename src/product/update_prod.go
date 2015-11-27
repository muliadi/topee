package product

import(
    "bytes"
    // "fmt"
    // "time"
    // "strconv"
    // "regexp"
    
    // "github.com/extemporalgenome/slug"
    // "gopkg.in/redis.v3"
    // "gopkg.in/mgo.v2/bson"
)

func UpdateProduct(input *ProductInput, current *ProductInput){
    
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
    
    //if catalog is changed
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
    
    //scan any phone number in description and insert to mongo for security team to use it later
    ScanPhoneNumber(input.ShortDesc, input.ProductId)
    
}


func RemoveFromCatalog(prod_id int64){
    buff := bytes.NewBufferString(`
        UPDATE ws_catalog_product 
        SET 
            status = 0,
            update_time = $1,
        WHERE
            product_id = $2,
        AND status = 1
    `)
    
    query := db.Rebind(buff.String())
    db.MustExec(
        query,
        prod_id,
        Now())
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
