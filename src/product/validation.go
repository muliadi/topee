package product

import(
    // "fmt"
    "sort"
    "regexp"
    // "net/http"
    // "log"
)

func ValidateInput(product *ProductInput) ([]Error){
    var errors []Error
    var error Error
    
    if CheckUserId(product.User_id) == false{
        error.Code      = "400"
        error.Message   = "user_id is invalid"
        errors = append(errors, error)
    }
    
    if len(product.Product_name) < 1{
        error.Code      = "400"
        error.Message   = "Product name must not empty"
        errors = append(errors, error)
    }
    if CheckProductName(product.Product_name) == false{
        error.Code      = "400"
        error.Message   = "Product name must contains alphanumeric"
        errors = append(errors, error)
    }
    
    if product.Shop_id == 0 {
        error.Code      = "400"
        error.Message   = "Shop id is required"
        errors = append(errors, error)
    } else if CheckShopId(product.Shop_id) == false {
        error.Code      = "400"
        error.Message   = "Shop id invalid"
        errors = append(errors, error)
    }
    
    if product.Child_cat_id == 0 {
        error.Code      = "400"
        error.Message   = "Category id is required"
        errors = append(errors, error)
    } else if CheckCategoryId(product.Child_cat_id) == false {
        error.Code      = "400"
        error.Message   = "Category id invalid"
        errors = append(errors, error)
    }
    
    if product.Min_order < 1 {
        error.Code      = "400"
        error.Message   = "Minimum order is required and must be greater than 0"
        errors = append(errors, error)
    }
    
    //set default price currency to IDR if not set
    if product.Price_currency == 0 {
        product.Price_currency = 1
    } else if product.Price_currency != 1 && product.Price_currency != 2{
        error.Code      = "400"
        error.Message   = "Price currency is not valid"
        errors = append(errors, error)
    }
    
    if product.Normal_price < 100 || product.Normal_price > 50000000 {
        error.Code      = "400"
        error.Message   = "Normal price must between 100 to 50.000.000"
        errors = append(errors, error)
    }
    
    if product.Weight.Numeric == 0{
        error.Code      = "400"
        error.Message   = "Weight is required"
        errors = append(errors, error)
    }
    
    // set default weight unit to gram
    if product.Weight.Unit == 0{
        product.Weight.Unit = 1
        
        if product.Weight.Numeric > 100000{
            error.Code      = "400"
            error.Message   = "Weight numeric must between 1 to 100.000"
            errors = append(errors, error)
        }
    } else if product.Weight.Unit != 1 && product.Weight.Unit != 2 {
        error.Code      = "400"
        error.Message   = "Weight unit value is not valid"
        errors = append(errors, error)
    }
    
    //max weight is 100kg
    if (product.Weight.Unit == 1 && product.Weight.Numeric > 100000) || (product.Weight.Unit == 2 && product.Weight.Numeric > 100){
        error.Code      = "400"
        error.Message   = "Max weight is 100 Kg"
        errors = append(errors, error)
    }
    
    if product.Insurance != 0 && product.Insurance !=1{
        error.Code      = "400"
        error.Message   = "must_insurance value is not valid"
        errors = append(errors, error)
    }
    
    if product.Add_to_etalase != 0 && product.Add_to_etalase !=1{
        error.Code      = "400"
        error.Message   = "add_to_etalase value is not valid"
        errors = append(errors, error)
    } else if product.Add_to_etalase == 1 && CheckEtalaseId(product.Etalase_id, product.Shop_id) == false{
        error.Code      = "400"
        error.Message   = "etalase_id is invalid"
        errors = append(errors, error)
    }
    
    //set product condition to 1 (new) if not set
    if product.Condition == 0 {
        product.Condition = 1
    } else if product.Condition != 1 && product.Condition != 2{
        error.Code      = "400"
        error.Message   = "condition value is invalid"
        errors = append(errors, error)
    }

    if len(product.Wholesale) > 0{
        var whole_bool bool
        whole_bool, product.Wholesale = ValidateWholesale(product.Wholesale, product.Normal_price)
        
        if whole_bool == false{
            error.Code      = "400"
            error.Message   = "wholesale is invalid"
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
