package product

import(
    "github.com/gin-gonic/gin"
    "github.com/extemporalgenome/slug"
    // "gopkg.in/redis.v3"
    // "fmt"
    "time"
    "strconv"
)

func Index(c *gin.Context) {
    rows, _ := db.Query(`SELECT product_id, shop_id, child_cat_id, product_name, short_desc FROM ws_product LIMIT 10`)
    defer rows.Close()
    
    var products []Product
    
    for rows.Next() {
        var prod Product
        
        err := rows.Scan(&prod.Product_id, &prod.Shop_id, &prod.Child_cat_id, &prod.Product_name, &prod.Short_desc)
        checkErr(err, "scan failed")
        
        products = append(products, prod)
        
    }
    c.JSON(200, products)
}

func Create(c *gin.Context){
    var input ProductInput
    c.Bind(&input)
    errors := ValidateInput(&input)
    
    if len(errors) > 0{
        c.JSON(400, errors)
    } else if len(errors) == 0 {
        
        //insert product into ws_product
        // product_id := InsertProduct(input)
        // CreateAlias(product_id, input.Product_name, input.Shop_id)
        
        c.JSON(200, input)
    }
}


// MAX POSITION FUNCTION - START
func GetMaxPosition(shop_id int64) int64{
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

func InsertProduct(product ProductInput) int64{
    //get max position in current shop
    current_pos := GetMaxPosition(product.Shop_id)
    current_pos++
    
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
        VALUES(
            $1,
            $2,
            $3,
            $4,
            $5,
            $6,
            $7,
            $8,
            $9,
            $10,
            $11,
            $12,
            $13,
            $14,
            $15,
            $16,
            $17
        )
        RETURNING product_id`,
        product.Shop_id,
        product.Child_cat_id,
        product.Product_name,
        product.Short_desc,
        product.Normal_price,
        product.Price_currency,
        product.Status,
        product.Weight.Numeric,
        product.Weight.Unit,
        product.Min_order,
        product.User_id,
        now,
        current_pos,
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
    err = cmgo.Insert(alias)
    checkErr(err, "Fail create alias in mongodb")
    
    //delete product alias in redis
    
}

func Now() string{
    time.LoadLocation("Asia/Jakarta")
    t := time.Now()
    now := t.Format("20060102150405")
    year := now[0:4]
    month := now[4:6]
    day := now[6:8]
    hour := now[8:10]
    min := now[10:12]
    sec := now[12:14]
    
    var timestamp string
    timestamp = year + "-" + month + "-" + day + " " + hour + ":" + min + ":" + sec
    return timestamp
}
