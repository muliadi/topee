package product

import(
    // "fmt"
    "bytes"
    "sort"
    "regexp"
    "strconv"
    "time"
    // "errors"
    // "net/http"
    // "log"
)

func ValidateInput(product *ProductInput) ([]Error){
    var errors []Error
    var error Error
    
    if CheckUserId(product.UserId) == false{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/user_id"
        error.Message           = "user_id is invalid"
        errors = append(errors, error)
    }
    
    if len(product.ProductName) < 1 || len(product.ProductName) > 70{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/name"
        error.Message           = "Product name must not empty and max 70 character"
        errors = append(errors, error)
    }
    if CheckProductName(product.ProductName) == false{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/name"
        error.Message           = "Product name must contains alphanumeric"
        errors = append(errors, error)
    }
    
    if product.ShopId == 0 {
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/shop_id"
        error.Message           = "Shop id is required"
        errors = append(errors, error)
    } else if CheckShopId(product.ShopId) == false {
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/shop_id"
        error.Message           = "Shop id invalid"
        errors = append(errors, error)
    }
    
    if product.ChildCatId == 0 {
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/category_id"
        error.Message           = "Category id is required"
        errors = append(errors, error)
    } else if CheckCategoryId(product.ChildCatId) == false {
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/category_id"
        error.Message           = "Category id invalid"
        errors = append(errors, error)
    }
    
    if product.MinOrder < 1 {
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/min_order"
        error.Message           = "Minimum order is required and must be greater than 0"
        errors = append(errors, error)
    }
    
    //set default price currency to IDR if not set
    if product.PriceCurrency == 0 {
        product.PriceCurrency = 1
    } else if product.PriceCurrency != 1 && product.PriceCurrency != 2{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/price_currency"
        error.Message           = "Price currency is not valid"
        errors = append(errors, error)
    }
    
    if product.NormalPrice < 100 || product.NormalPrice > 50000000 {
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/normal_price"
        error.Message           = "Normal price must between 100 to 50.000.000"
        errors = append(errors, error)
    }
    
    if product.Weight.Numeric == 0{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/weight/numeric"
        error.Message           = "Weight is required"
        errors = append(errors, error)
    }
    
    // set default weight unit to gram
    if product.Weight.Unit == 0{
        product.Weight.Unit = 1
        
        if product.Weight.Numeric > 100000{
            // error.Code              = "400"
            error.Source.Pointer    = "/data/attributes/weight/numeric"
            error.Message           = "Weight numeric must between 1 to 100.000"
            errors = append(errors, error)
        }
    } else if product.Weight.Unit != 1 && product.Weight.Unit != 2 {
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/weight/unit"
        error.Message           = "Weight unit value is not valid"
        errors = append(errors, error)
    }
    
    //max weight is 100kg
    if (product.Weight.Unit == 1 && product.Weight.Numeric > 100000) || (product.Weight.Unit == 2 && product.Weight.Numeric > 100){
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/weight/numeric"
        error.Message           = "Max weight is 100 Kg"
        errors = append(errors, error)
    }
    
    //weight must not 0
    if product.Weight.Unit == 0{
        error.Source.Pointer    = "/data/attributes/weight/numeric"
        error.Message           = "weight unit is required and must larger than 0"
        errors = append(errors, error)
    }
    
    if product.Insurance != 0 && product.Insurance !=1{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/must_insurance"
        error.Message           = "must_insurance value is not valid"
        errors = append(errors, error)
    }
    
    //etalase
    if product.AddToEtalase != 0 && product.AddToEtalase !=1{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/add_to_etalase"
        error.Message           = "add_to_etalase value is not valid"
        errors = append(errors, error)
    } else if product.AddToEtalase == 1 && CheckEtalaseId(product.EtalaseId, product.ShopId) == false{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/etalase_id"
        error.Message           = "etalase_id is invalid"
        errors = append(errors, error)
    }
    
    //catalog
    if product.AddToCatalog != 0 && product.AddToCatalog !=1{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/add_to_catalog"
        error.Message           = "add_to_catalog value is not valid"
        errors = append(errors, error)
    } else if product.AddToCatalog == 1 && CheckCatalogId(product.CatalogId, product.ChildCatId) == false{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/catalog_id"
        error.Message           = "catalog_id is invalid"
        errors = append(errors, error)
    }
    
    //validasi product status
    // if product.Status != 0 && product.Status !=1 && product.Status !=2 && product.Status !=3 && product.Status !=-1 && product.Status !=-2{
    //     // error.Code              = "400"
    //     error.Source.Pointer    = "/data/attributes/product_status"
    //     error.Message           = "product_status value is not valid"
    //     errors = append(errors, error)
    // }
    
    //validasi product returnable
    if product.Returnable != 0 && product.Returnable !=1{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/returnable"
        error.Message           = "returnable value is not valid"
        errors = append(errors, error)
    }
    
    //set product condition to 1 (new) if not set
    if product.Condition == 0 {
        product.Condition = 1
    } else if product.Condition != 1 && product.Condition != 2{
        // error.Code              = "400"
        error.Source.Pointer    = "/data/attributes/condition"
        error.Message           = "condition value is invalid"
        errors = append(errors, error)
    }

    if len(product.Wholesale) > 0{
        var whole_bool bool
        whole_bool, product.Wholesale = ValidateWholesale(product.Wholesale, product.NormalPrice)
        
        if whole_bool == false{
            // error.Code              = "400"
            error.Source.Pointer    = "/data/attributes/wholesale"
            error.Message           = "wholesale is invalid"
            errors = append(errors, error)
        }
    }
    
    return errors
}

func ValidateWholesale(wsprices WsPrices, normal_price int64) (bool, WsPrices){
    sort.Sort(wsprices)
    
    var price_before = normal_price
    var qty_before int64 = 0
    var error int = 0
    for _, ws := range wsprices {
        if ws.Min >= ws.Max {
            error = 1
        } 
        if ws.Min <= qty_before {
            error = 1
        }
        if ws.Price >= price_before{
            error = 1
        }
        
        price_before = ws.Price
        qty_before = ws.Max
    }
    
    if error == 0 {
        return true, wsprices
    } else {
        return false, wsprices
    }
}

func CheckCategoryId(id int64) bool{
    var count int
    row := db.QueryRow("SELECT count(d_id) FROM ws_department WHERE d_id=$1 OR parent = $2", id, id)
    row.Scan(&count)
    if count == 1{
        return true
    } else {
        return false
    }
}

func CheckCatalogId(cat_id int64, dep_id int64) bool{
    var count int
    row := db.QueryRow("SELECT count(ctg_id) FROM ws_catalog WHERE ctg_id=$1 AND d_id=$2 AND status=1", cat_id, dep_id)
    row.Scan(&count)
    if count == 1{
        return true
    } else {
        return false
    }
}

func CheckShopId(id int64) bool{
    var count int
    row := db.QueryRow("SELECT count(shop_id) FROM ws_shop WHERE shop_id = $1", id) 
    row.Scan(&count)
    if count == 1{
        return true
    } else {
        return false
    }
}

func CheckEtalaseId(etalase_id int64, shop_id int64) bool{
    var count int
    row := db.QueryRow("SELECT count(menu_id) FROM ws_menu WHERE menu_id = $1 AND shop_id = $2", etalase_id, shop_id)
    row.Scan(&count)
    if count == 1{
        return true
    } else {
        return false
    }
}

func CheckUserId(user_id int64) bool{
    var count int
    row := db.QueryRow("SELECT count(user_id) FROM ws_user WHERE user_id = $1", user_id)
    row.Scan(&count)
    if count == 1{
        return true
    } else {
        return false
    }
}

func CheckProductName(name string) bool{
    match, _ := regexp.MatchString("[a-zA-Z0-9]", name)
    return match
}

func CheckProductId(prod_id int64) (bool, error){
    var res int64
    buff := bytes.NewBufferString(`
        SELECT product_id FROM ws_product
        WHERE product_id = $1
    `)
    
    query := db.Rebind(buff.String())
    err := db.QueryRow(query, prod_id).Scan(&res)
    // fmt.Println(res)
    if res == prod_id {
        return true, err
    } else {
        return false, err
    }
}

func GetProductById(prod_id int64) (ProductInput, error){
    var prd ProductInput
    buff := bytes.NewBufferString(`
        SELECT 
            product_id,
            product_name,
            shop_id,
            child_cat_id,
            short_desc,
            min_order,
            price_currency,
            normal_price,
            status,
            position,
            must_insurance,
            condition,
            weight_unit,
            weight,
            last_update_price
        FROM ws_product
        WHERE product_id = $1 LIMIT 1
    `)
    query := db.Rebind(buff.String())
    row := db.QueryRow(query, prod_id)
    
    var normal_price string
    var lastprc time.Time
    
    err := row.Scan(
        &prd.ProductId,
        &prd.ProductName,
        &prd.ShopId,
        &prd.ChildCatId,
        &prd.ShortDesc,
        &prd.MinOrder,
        &prd.PriceCurrency,
        &normal_price,
        &prd.Status,
        &prd.Position,
        &prd.Insurance,
        &prd.Condition,
        &prd.Weight.Unit,
        &prd.Weight.Numeric,
        &lastprc)
    
    prd.LastPrcUpdate = lastprc.String()
    
    prd.NormalPrice, _ = strconv.ParseInt(normal_price, 10, 64)
    
    if err != nil{
        return ProductInput{}, err
    } else {
        return prd, nil
    }
}


func ValidateUpdate(product *ProductInput) (ProductInput, []Error, error){
    var errors []Error
    var error Error
    var current_prod ProductInput
    
    var changes int64 = 0
    
    //validate product id
    res, err := CheckProductId(product.ProductId)
    if err != nil {
        // return nil , err
    }
    if res == false {
        error.Source.Pointer    = "/data/id"
        error.Message           = "id is invalid"
        errors = append(errors, error)
    } else if res == true{
        //if product id is valid
        //get current product data
        
        current_prod, err = GetProductById(product.ProductId)
        if err != nil{
            return ProductInput{}, nil, err
        }
        
        //if category is updated
        if product.ChildCatId != 0{
            if CheckCategoryId(product.ChildCatId) == false{
                error.Source.Pointer    = "/data/attributes/category_id"
                error.Message           = "Category id invalid"
                errors = append(errors, error)
            }
            changes++
        }
        
        //if price currency is updated
        if product.PriceCurrency != 0 && product.PriceCurrency != 1 && product.PriceCurrency != 2{
            error.Source.Pointer    = "/data/attributes/price_currency"
            error.Message           = "Price currency is not valid"
            errors = append(errors, error)
        }
        
        //if normal price is updated
        if product.NormalPrice != 0 && (product.NormalPrice < 100 || product.NormalPrice > 50000000){
            error.Source.Pointer    = "/data/attributes/normal_price"
            error.Message           = "Normal price must between 100 to 50.000.000"
            errors = append(errors, error)
        }
        
        
        //if weight is updated
        if product.Weight.Unit != 0 && product.Weight.Numeric != 0{
            //check weight unit
            if product.Weight.Unit != 1 && product.Weight.Unit != 2 {
                // error.Code              = "400"
                error.Source.Pointer    = "/data/attributes/weight/unit"
                error.Message           = "Weight unit value is not valid"
                errors = append(errors, error)
            }
            
            //max weight is 100kg
            if (product.Weight.Unit == 1 && product.Weight.Numeric > 100000) || (product.Weight.Unit == 2 && product.Weight.Numeric > 100){
                // error.Code              = "400"
                error.Source.Pointer    = "/data/attributes/weight/numeric"
                error.Message           = "Max weight is 100 Kg"
                errors = append(errors, error)
            }
        }
        if product.Weight.Unit != 0 && product.Weight.Numeric == 0 {
            error.Source.Pointer    = "/data/attributes/weight/numeric"
            error.Message           = "weight numeric is required"
            errors = append(errors, error)
        }
        if product.Weight.Unit == 0 && product.Weight.Numeric != 0 {
            error.Source.Pointer    = "/data/attributes/weight/unit"
            error.Message           = "weight unit is required"
            errors = append(errors, error)
        }
        
        //if must insurance is updated
        if product.Insurance != -1 && product.Insurance != 0 && product.Insurance != 1{
            error.Source.Pointer    = "/data/attributes/must_insurance"
            error.Message           = "must_insurance value is not valid"
            errors = append(errors, error)
        }
        
        //if etalase is updated
        if product.AddToEtalase != -1 && product.EtalaseId != 0{
            if product.AddToEtalase != 0 && product.AddToEtalase !=1{
                // error.Code              = "400"
                error.Source.Pointer    = "/data/attributes/add_to_etalase"
                error.Message           = "add_to_etalase value is not valid"
                errors = append(errors, error)
            } else if product.AddToEtalase == 1 && CheckEtalaseId(product.EtalaseId, current_prod.ShopId) == false{
                // error.Code              = "400"
                error.Source.Pointer    = "/data/attributes/etalase_id"
                error.Message           = "etalase_id is invalid"
                errors = append(errors, error)
            }
        }
        if product.AddToEtalase != -1 && product.EtalaseId == 0{
            error.Source.Pointer    = "/data/attributes/etalase_id"
            error.Message           = "etalase id is required"
            errors = append(errors, error)
        }
        if product.AddToEtalase == -1 && product.EtalaseId != 0{
            error.Source.Pointer    = "/data/attributes/add_to_etalase"
            error.Message           = "wrong value on add_to_etalase"
            errors = append(errors, error)
        }
        
        //if catalog is updated
        if product.AddToCatalog != -1 && product.CatalogId != 0{
            if product.AddToCatalog != 0 && product.AddToCatalog !=1{
                // error.Code              = "400"
                error.Source.Pointer    = "/data/attributes/add_to_catalog"
                error.Message           = "add_to_catalog value is not valid"
                errors = append(errors, error)
            } else if product.AddToCatalog == 1{
                if product.ChildCatId == 0 && CheckCatalogId(product.CatalogId, current_prod.ChildCatId) == false{
                    // error.Code              = "400"
                    error.Source.Pointer    = "/data/attributes/Catalog_id"
                    error.Message           = "Catalog_id is invalid"
                    errors = append(errors, error)
                } else if product.ChildCatId != 0 && CheckCatalogId(product.CatalogId, product.ChildCatId) == false {
                    // error.Code              = "400"
                    error.Source.Pointer    = "/data/attributes/Catalog_id"
                    error.Message           = "Catalog_id is invalid"
                    errors = append(errors, error)
                }
            }
        }
        if product.AddToCatalog != -1 && product.CatalogId == 0{
            error.Source.Pointer    = "/data/attributes/Catalog_id"
            error.Message           = "Catalog id is required"
            errors = append(errors, error)
        }
        if product.AddToCatalog == -1 && product.CatalogId != 0{
            error.Source.Pointer    = "/data/attributes/add_to_Catalog"
            error.Message           = "wrong value on add_to_Catalog"
            errors = append(errors, error)
        }
        
        //if condition is updated
        if product.Condition != 0 && product.Condition != 1 && product.Condition != 2{
            // error.Code              = "400"
            error.Source.Pointer    = "/data/attributes/condition"
            error.Message           = "condition value is invalid"
            errors = append(errors, error)
        }
        
        //if returnable is updated
        if product.Returnable != -1 && product.Returnable != 0 && product.Returnable != 1{
            error.Source.Pointer    = "/data/attributes/returnable"
            error.Message           = "returnable value is not valid"
            errors = append(errors, error)
        }
        
        if len(product.Wholesale) > 0{
            var whole_bool bool
            var base_price int64
            if product.NormalPrice == 0{
                base_price = current_prod.NormalPrice
            } else {
                base_price = product.NormalPrice
            }
            whole_bool, product.Wholesale = ValidateWholesale(product.Wholesale, base_price)
            
            if whole_bool == false{
                // error.Code              = "400"
                error.Source.Pointer    = "/data/attributes/wholesale"
                error.Message           = "wholesale is invalid"
                errors = append(errors, error)
            }
        }
    }
    
    return current_prod, errors, nil
}
