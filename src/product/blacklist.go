package product

import(
    // "fmt"
    // "sort"
    "regexp"
    // "log"
    "bytes"
    "strings"
)

func CheckBlacklist(text string, blacklist_type int) (bool, string){
    
    //get blacklist list
    blacklists := BlacklistWords(blacklist_type)
    var words_slice []string
    for _, blist := range blacklists {
        words_slice = append(words_slice, blist.Value)
    }
    
    words_regex := strings.Join(words_slice, "|")
    
    re := regexp.MustCompile(words_regex)
    tmatch := re.FindStringSubmatch(text)
    
    if len(tmatch) > 0{
        return true, ""
    } else {
        match_words := strings.Join(tmatch, ",")
        return false, match_words
    }
}

func CheckBajakan(harga int64, cat_id int64) bool{
    categories := map[int64]bool{
        17      : true,
        1517    : true,
        19      : true,
        20      : true,
    }
    
    var max_harga int64 = 10000
    
    if harga <= max_harga && categories[cat_id] == true {
        return true
    } else {
        return false
    }
}

func BlacklistWords(blacklist_type int) []Blacklist{
    buff := bytes.NewBufferString(`
        SELECT
            blacklist_id,
            blacklist_value,
            blacklist_type,
            blacklist_reason,
            blacklist_status
        FROM
            ws_blacklist
        WHERE
            blacklist_type = $1 AND
            blacklist_status = 1
    `)
    
    query := db.Rebind(buff.String())
    blacklists := []Blacklist{}
    err := db.Select(&blacklists, query, blacklist_type)
    checkErr(err, "Fail to get blacklist from ws_blacklist")
    
   return blacklists
}

func UpsertPendingReason(product *ProductInput){
    buff := bytes.NewBufferString(`
        SELECT
            'x'
        FROM
            ws_product_pending_reason
        WHERE
            product_id = $1
    `)
    query := db_product.Rebind(buff.String())
    
    row := db_product.QueryRow(query, product.ProductId)
    var res string
    row.Scan(&res)
    
    if res == "x" {
        //if exist update
        buff2 := bytes.NewBufferString(`
            UPDATE ws_product_pending_reason SET
                pending_type    = $2,
                pending_reason  = $3,
                update_time     = $4
            WHERE
                product_id = $1
        `)
        query = db_product.Rebind(buff2.String())
        
    } else {
        //else insert
        buff2 := bytes.NewBufferString(`
            INSERT INTO ws_product_pending_reason (
                product_id,
                pending_type,
                pending_reason,
                create_time
            ) VALUES ($1, $2, $3, $4)
        `)
        query = db_product.Rebind(buff2.String())
    }
    db_product.MustExec(query, product.ProductId, product.PendingStatus, product.PendingReason, Now())
}
