package product

import(
    // "gopkg.in/mgo.v2/bson"
)

//==============================================================================
//  STRUCT FOR CONFIG
//==============================================================================
type Config struct {
    Database    map[string]string
    Mongo       map[string]string
    Redis       map[string]Redis
    Port        int
}

type Redis struct{
    Host        string
    Port        string
}



//==============================================================================
//  STRUCT FOR JSON
//==============================================================================
type Product struct {
    Product_id      int64           `json:"product_id"`
    Shop_id         int64           `json:"shop_id"`
    Child_cat_id    string          `json:"child_cat_id"`
    Product_name    string          `json:"product_name"`
    Short_desc      string          `json:"short_desc"`
}

type ProductInput struct{
    Product_id      int64           `json:"product_id"`
    User_id         int64           `json:"user_id"`
    Product_name    string          `json:"name"`
    Shop_id         int64           `json:"shop_id"`
    Child_cat_id    int64           `json:"category_id"`
    Short_desc      string          `json:"description"`
    Min_order       int64           `json:"min_order"`
    Price_currency  int64           `json:"price_currency"` // 1=IDR 2=USD
    Normal_price    int64           `json:"normal_price"`
    Status          int64           `json:"product_status"` // 0=deleted, 1=active, 2=best, 3=warehouse, -1=pending, -2=banned
    Position        int64           `json:"position"`
    Weight          ProductWeight   `json:"weight"`
    Insurance       int64           `json:"must_insurance"` // 0=no 1=yes
    Add_to_etalase  int64           `json:"add_to_etalase"` // 0=no 1=yes
    Etalase_id      int64           `json:"etalase_id"`
    Condition       int64           `json:"condition"`      // 1=baru 2=bekas
    Returnable      int64           `json:"returnable"`     // 0=no 1=yes
    Wholesale       WsPrices        `json:"wholesale"`
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
    ProductId       int64           `bson:"product_id"`
    ProductKey      string          `bson:"product_key"`
    ShopId          int64           `bson:"shop_id"`
}

type WholesaleMongo struct{
    UpdateTime      int64           `bson:"update_time"`
    ProductId       int64           `bson:"product_id"`
    UpdateBy        int64           `bson:"update_by"`
    QtyMin1         int64           `bson:"qty_min_1, omitempty"`
    QtyMax1         int64           `bson:"qty_max_1, omitempty"`
    PrdPrc1         int64           `bson:"prd_prc_1, omitempty"`
    QtyMin2         int64           `bson:"qty_min_2, omitempty"`
    QtyMax2         int64           `bson:"qty_max_2, omitempty"`
    PrdPrc2         int64           `bson:"prd_prc_2, omitempty"`
    QtyMin3         int64           `bson:"qty_min_3, omitempty"`
    QtyMax3         int64           `bson:"qty_max_3, omitempty"`
    PrdPrc3         int64           `bson:"prd_prc_3, omitempty"`
    QtyMin4         int64           `bson:"qty_min_4, omitempty"`
    QtyMax4         int64           `bson:"qty_max_4, omitempty"`
    PrdPrc4         int64           `bson:"prd_prc_4, omitempty"`
    QtyMin5         int64           `bson:"qty_min_5, omitempty"`
    QtyMax5         int64           `bson:"qty_max_5, omitempty"`
    PrdPrc5         int64           `bson:"prd_prc_5, omitempty"`
}

type ProductInfoMongo struct{
    ProductId       int64           `bson:"product_id"`
    ShopId          int64           `bson:"shop_id"`
    Returnable      int64           `bson:"returnable"`
}
