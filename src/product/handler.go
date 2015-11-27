package product

import(
    "fmt"
    "net/http"
    "encoding/json"
    "strconv"
    
    "github.com/julienschmidt/httprouter"
    au "github.com/ruizu/api-utils"
)


func Ping(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    w.Write([]byte("Pong"))
    
    CheckProductId(213301)
    // CheckBlacklist("ini adalah produk dengan kata asdasdasd didalamnya", 4)
    // var deskripsi string = `
    //     testing deskripsi produk dengan nomer HP seperti 
    //     085851601112 +628565160110
    //      atau 085851601112
    // `
    // ScanPhoneNumber(deskripsi, 1234)
    //TODO: Health check, as in https://github.com/tokopedia/search-microservice/blob/master/app.go#L33
}

func CreateHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params){
    var input ProductPostRaw
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&input)
    checkErr(err, "failed to decode JSON input")
    
    if input.Data.Type == "products" {
        product := input.Data.Attributes
        errors := ValidateInput(&product)
        fmt.Println(product)
        
        if len(errors) > 0{
            res := new(ResErrors)
            res.Links = new(au.ResponseLink)
            res.Links.Self = r.URL.String()
            res.Data = errors
            au.WriteResponse(w, res, http.StatusBadRequest)
        } else if len(errors) == 0 {
            
            CreateProduct(&product)
           
            res := new(ResCreateProd)
            res.Links = new(au.ResponseLink)
            res.Links.Self = GetProductUri(&product)
            res.Data = product
           
            au.WriteResponse(w, res, http.StatusCreated)
        }
    } else {
        var errors []Error
        var error Error
        
        error.Source.Pointer    = "/data/type"
        error.Message           = "unknown type value"
        errors = append(errors, error)
        
        res := new(ResErrors)
        res.Links = new(au.ResponseLink)
        res.Links.Self = r.URL.String()
        res.Data = errors
        au.WriteResponse(w, res, http.StatusBadRequest)
    }
}

func UpdateHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params){
    input := ProductPostRaw{
        Data    : RawData{
            Type        : "",
            Id          : 0,
            Attributes  : ProductInput{
                Returnable      : -1,
                Insurance       : -1,
                AddToEtalase    : -1,
                AddToCatalog    : -1,
            },
        },
    }
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&input)
    checkErr(err, "failed to decode JSON input")
    param_id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
    
    if input.Data.Type != "products" {
        var errors []Error
        var error Error
        
        error.Source.Pointer    = "/data/type"
        error.Message           = "unknown type value"
        errors = append(errors, error)
        
        res := new(ResErrors)
        res.Links = new(au.ResponseLink)
        res.Links.Self = r.URL.String()
        res.Data = errors
        au.WriteResponse(w, res, http.StatusBadRequest)
    } else if input.Data.Id != param_id{
        var errors []Error
        var error Error
        
        error.Source.Pointer    = "/data/id"
        error.Message           = "wrong id"
        errors = append(errors, error)
        
        res := new(ResErrors)
        res.Links = new(au.ResponseLink)
        res.Links.Self = r.URL.String()
        res.Data = errors
        au.WriteResponse(w, res, http.StatusBadRequest)
    } else {
        product := input.Data.Attributes
        product.ProductId = input.Data.Id
        current_prod, errors, err := ValidateUpdate(&product)
        
        if err != nil{
            return
        }
        if len(errors) > 0{
            res := new(ResErrors)
            res.Links = new(au.ResponseLink)
            res.Links.Self = r.URL.String()
            res.Data = errors
            au.WriteResponse(w, res, http.StatusBadRequest)
        } else if len(errors) == 0{
            UpdateProduct(&product, &current_prod)
        }
    }
}
