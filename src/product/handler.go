package product

import(
   
    "net/http"
    "encoding/json"
    
    "github.com/julienschmidt/httprouter"
    au "github.com/ruizu/api-utils"
)


func Ping(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    w.Write([]byte("Pong"))
    CheckBlacklist("ini adalah produk dengan kata asdasdasd didalamnya", 4)
    //TODO: Health check, as in https://github.com/tokopedia/search-microservice/blob/master/app.go#L33
}

func IndexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    rows, _ := db.Query(`SELECT product_id, shop_id, child_cat_id, product_name, short_desc FROM ws_product LIMIT 10`)
    defer rows.Close()
    
    var products []Product
    
    res := new(ResIndexProd)
    res.Links = new(au.ResponseLink)
    res.Links.Self = r.URL.String()
    
    for rows.Next() {
        var prod Product
        
        err := rows.Scan(&prod.ProductId, &prod.ShopId, &prod.ChildCatId, &prod.ProductName, &prod.ShortDesc)
        checkErr(err, "scan failed")
        
        products = append(products, prod)
    }
    res.Data = products
    au.WriteResponse(w, res, http.StatusOK)
}

func CreateHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params){
    var input ProductInput
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&input)
    checkErr(err, "failed to decode JSON input")
    errors := ValidateInput(&input)
    
    
    if len(errors) > 0{
        res := new(ResErrors)
        res.Links = new(au.ResponseLink)
        res.Links.Self = r.URL.String()
        res.Data = errors
        au.WriteResponse(w, res, http.StatusBadRequest)
    } else if len(errors) == 0 {
        
        CreateProduct(&input)
       
        res := new(ResCreateProd)
        res.Links = new(au.ResponseLink)
        res.Links.Self = r.URL.String()
        res.Data = input
       
        au.WriteResponse(w, res, http.StatusCreated)
    }
}
