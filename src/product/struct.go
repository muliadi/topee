package product

import(
    "gopkg.in/mgo.v2/bson"
    au "github.com/ruizu/api-utils"
)

//==============================================================================
//  STRUCT FOR RESPONSE
//==============================================================================
type ResCreateProd struct{
    au.Response
    Data            ProductInput    `json:"data,omitempty"`
}

type ResIndexProd struct{
    au.Response
    Data            []Product       `json:"data,omitempty"`
}

type ResErrors struct{
    au.Response
    Data            []Error         `json:"errors,omitempty"`
}

//==============================================================================
//  STRUCT FOR JSON DATA
//==============================================================================
type Product struct {
    ProductId       int64           `json:"product_id"`
    ShopId          int64           `json:"shop_id"`
    ChildCatId      string          `json:"child_cat_id"`
    ProductName     string          `json:"product_name"`
    ShortDesc       string          `json:"short_desc"`
}

type ProductInput struct{
    ProductId       int64           `json:"product_id"`
    UserId          int64           `json:"user_id"`
    ProductName     string          `json:"name"`
    ShopId          int64           `json:"shop_id"`
    ChildCatId      int64           `json:"category_id"`
    CatalogId       int64           `json:"catalog_id"`
    ShortDesc       string          `json:"description"`
    MinOrder        int64           `json:"min_order"`
    PriceCurrency   int64           `json:"price_currency"` // 1=IDR 2=USD
    NormalPrice     int64           `json:"normal_price"`
    Status          int64           `json:"product_status"` // 0=deleted, 1=active, 2=best, 3=warehouse, -1=pending, -2=banned
    Position        int64           `json:"position"`
    Weight          ProductWeight   `json:"weight"`
    Insurance       int64           `json:"must_insurance"` // 0=no 1=yes
    AddToEtalase    int64           `json:"add_to_etalase"` // 0=no 1=yes
    EtalaseId       int64           `json:"etalase_id"`
    Condition       int64           `json:"condition"`      // 1=baru 2=bekas
    Returnable      int64           `json:"returnable"`     // 0=no 1=yes
    Wholesale       WsPrices        `json:"wholesale"`
    
    PendingReason   string          `json:-`
    PendingStatus   int64           `json:-`
}

type ProductWeight struct{
    Unit            int64           `json:"unit"`           // 1=gram 2=Kg
    Numeric         float64         `json:"numeric"`
}

type WholesalePrice struct{
    Min             int64           `json:"min_count"`
    Max             int64           `json:"max_count"`
    Price           int64           `json:"price"`
}

type WsPrices []WholesalePrice
func (slice WsPrices) Len() int {
    return len(slice)
}
func (slice WsPrices) Less(i, j int) bool {
    return slice[i].Min < slice[j].Min;
}
func (slice WsPrices) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

type Error struct{
    Code            string          `json:"code"`
    Message         string          `json:"message"`
}




//==============================================================================
//  STRUCT FOR MONGODB
//==============================================================================
type ProductAlias struct{
    Id              bson.ObjectId   `bson:"_id,omitempty"`
    ProductId       int64           `bson:"product_id"`
    ProductKey      string          `bson:"product_key"`
    ShopId          int64           `bson:"shop_id"`
}

type WholesaleMongo struct{
    Id              bson.ObjectId   `bson:"_id,omitempty"`
    UpdateTime      int64           `bson:"update_time"`
    ProductId       int64           `bson:"product_id"`
    UpdateBy        int64           `bson:"update_by"`
    QtyMin1         int64           `bson:"qty_min_1,omitempty"`
    QtyMax1         int64           `bson:"qty_max_1,omitempty"`
    PrdPrc1         int64           `bson:"prd_prc_1,omitempty"`
    QtyMin2         int64           `bson:"qty_min_2,omitempty"`
    QtyMax2         int64           `bson:"qty_max_2,omitempty"`
    PrdPrc2         int64           `bson:"prd_prc_2,omitempty"`
    QtyMin3         int64           `bson:"qty_min_3,omitempty"`
    QtyMax3         int64           `bson:"qty_max_3,omitempty"`
    PrdPrc3         int64           `bson:"prd_prc_3,omitempty"`
    QtyMin4         int64           `bson:"qty_min_4,omitempty"`
    QtyMax4         int64           `bson:"qty_max_4,omitempty"`
    PrdPrc4         int64           `bson:"prd_prc_4,omitempty"`
    QtyMin5         int64           `bson:"qty_min_5,omitempty"`
    QtyMax5         int64           `bson:"qty_max_5,omitempty"`
    PrdPrc5         int64           `bson:"prd_prc_5,omitempty"`
}

type ProductInfoMongo struct{
    Id              bson.ObjectId   `bson:"_id,omitempty"`
    ProductId       int64           `bson:"product_id"`
    ShopId          int64           `bson:"shop_id"`
    Returnable      int64           `bson:"returnable"`
}

type ProductListMongo struct{
    Id              bson.ObjectId   `bson:"_id,omitempty"`
    ChildCatId      int64           `bson:"child_cat_id"`
    Condition       int64           `bson:"condition"`
    CountReview     int64           `bson:"count_review"`
    CountTalk       int64           `bson:"count_talk"`
    CountView       int64           `bson:"count_view"`
    CreateTime      int64           `bson:"create_time"`
    CatalogId       int64           `bson:"ctg_id"`
    IsVerified      int64           `bson:"is_verified"`
    ItemSold        int64           `bson:"item_sold"`
    MenuId          int64           `bson:"menu_id"`
    NormalPrice     int64           `bson:"normal_price"`
    PictureId       int64           `bson:"picture_id"`
    Position        int64           `bson:"position"`
    PriceCurrency   int64           `bson:"price_currency"`
    ProductId       int64           `bson:"product_id"`
    ProductName     string          `bson:"product_name"`
    RupiahPrice     int64           `bson:"rupiah_price"`
    ShopId          int64           `bson:"shop_id"`
    Status          int64           `bson:"status"`
    TempCategory    int64           `bson:"temp_category"`
    UpdateTime      int64           `bson:"update_time"`
}


//==============================================================================
//  STRUCT FOR BLACKLIST
//==============================================================================
type Blacklist struct{
    Id              int64           `db:"blacklist_id"`
    Value           string          `db:"blacklist_value"`
    Type            int64           `db:"blacklist_type"`
    Reason          string          `db:"blacklist_reason"`
    Status          int64           `db:"blacklist_status"`
}
