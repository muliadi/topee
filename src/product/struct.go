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

type ResErrors struct{
    au.Response
    Data            []Error         `json:"errors,omitempty"`
}

//==============================================================================
//  STRUCT FOR JSON DATA
//==============================================================================
type ProductPostRaw struct{
    Data            RawData         `json:"data"`
}

type RawData struct{
    Type            string          `json:"type"`
    Id              int64           `json:"id"`
    Attributes      ProductInput    `json:"attributes"`
}

type ProductInput struct{
    ProductId       int64           `json:"product_id"      db:"product_id"`
    UserId          int64           `json:"user_id"         db:"-"`
    ProductName     string          `json:"name"            db:"product_name"`
    ShopId          int64           `json:"shop_id"         db:"shop_id"`
    ChildCatId      int64           `json:"category_id"     db:"child_cat_id"`
    ShortDesc       string          `json:"description"     db:"short_desc"`
    MinOrder        int64           `json:"min_order"       db:"min_order"`
    PriceCurrency   int64           `json:"price_currency"  db:"price_currency"` // 1=IDR 2=USD
    NormalPrice     int64           `json:"normal_price"    db:"normal_price"`
    
    ProductStatus   string          `json:"product_status"  db:"-"`
    Status          int64           `json:"-"               db:"status"` // 0=deleted, 1=active, 2=best, 3=warehouse, -1=pending, -2=banned
    
    Position        int64           `json:"position"        db:"position"`
    Weight          ProductWeight   `json:"weight"          db:"-"`
    Insurance       int64           `json:"must_insurance"  db:"must_insurance"` // 0=no 1=yes
    AddToEtalase    int64           `json:"add_to_etalase"  db:"-"` // 0=no 1=yes
    EtalaseId       int64           `json:"etalase_id"      db:"-"`
    AddToCatalog    int64           `json:"add_to_catalog"  db:"-"` // 0=no 1=yes
    CatalogId       int64           `json:"catalog_id"      db:"-"`
    Condition       int64           `json:"condition"       db:"condition"`      // 1=baru 2=bekas
    Returnable      int64           `json:"returnable"      db:"-"`     // 0=no 1=yes
    Wholesale       WsPrices        `json:"wholesale"       db:"-"`
    Alias           string          `json:"product_alias"   db:"-"`
    
    PendingReason   string          `json:"-"               db:"-"`
    PendingStatus   int64           `json:"-"               db:"-"`
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
    Code            string          `json:"status,omitempty"`
    Source          ErrorSource     `json:"source"`
    Message         string          `json:"detail"`
}

type ErrorSource struct{
    Pointer         string          `json:"pointer"`
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
    ChildCatId      int64           `bson:"child_cat_id,omitempty"`
    Condition       int64           `bson:"condition,omitempty"`
    CountReview     int64           `bson:"count_review,omitempty"`
    CountTalk       int64           `bson:"count_talk,omitempty"`
    CountView       int64           `bson:"count_view,omitempty"`
    CreateTime      int64           `bson:"create_time,omitempty"`
    CatalogId       int64           `bson:"ctg_id,omitempty"`
    IsVerified      int64           `bson:"is_verified,omitempty"`
    ItemSold        int64           `bson:"item_sold,omitempty"`
    MenuId          int64           `bson:"menu_id,omitempty"`
    NormalPrice     int64           `bson:"normal_price,omitempty"`
    PictureId       int64           `bson:"picture_id,omitempty"`
    Position        int64           `bson:"position,omitempty"`
    PriceCurrency   int64           `bson:"price_currency,omitempty"`
    ProductId       int64           `bson:"product_id,omitempty"`
    ProductName     string          `bson:"product_name,omitempty"`
    RupiahPrice     int64           `bson:"rupiah_price,omitempty"`
    ShopId          int64           `bson:"shop_id,omitempty"`
    Status          int64           `bson:"status,omitempty"`
    TempCategory    int64           `bson:"temp_category,omitempty"`
    UpdateTime      int64           `bson:"update_time,omitempty"`
}

type UserLog struct{
    Id              bson.ObjectId   `bson:"_id,omitempty"`
    CreateTime      int64           `bson:"create_time"`
    Action          int64           `bson:"action"`
    ProductId       int64           `bson:"id"`
    UserId          int64           `bson:"user_id"`
    IpAddress       string          `bson:"ip_address"`
    Device          int64           `bson:"device"`
}

type ProductHistory struct{
    Id              bson.ObjectId   `bson:"_id,omitempty"`
    PriceCurrency   int64           `bson:"price_currency"`
    CreateBy        int64           `bson:"create_by"`
    Status          int64           `bson:"status"`
    CreateTime      int64           `bson:"create_time"`
    MinOrder        int64           `bson:"min_order"`
    MustInsurance   string          `bson:"must_insurance,omitempty"`
    WeightUnit      int64           `bson:"weight_unit"`
    ProductId       int64           `bson:"product_id"`
    NormalPrice     int64           `bson:"normal_price"`
    Weight          float64         `bson:"weight"`
    ShortDesc       string          `bson:"short_desc"`
    ChildCatId      int64           `bson:"child_cat_id"`
    ProductName     string          `bson:"product_name"`
    Position        int64           `bson:"position,omitempty"`
    Condition       int64           `bson:"condition,omitempty"`
    Returnable      int64           `bson:"returnable,omitempty"`
}

type PhoneNumber struct{
    Id              bson.ObjectId   `bson:"_id,omitempty"`
    PhoneNumber     string          `bson:"phone_number"`
    ProductId       int64           `bson:"product_id"`
    Description     string          `bson:"description"`
    Status          int64           `bson:"status"`
    CreateTime      int64           `bson:"create_time"`
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
