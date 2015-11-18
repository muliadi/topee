package product

type Product struct {
    Product_id      int64   `json:"product_id"`
    Shop_id         int64   `json:"shop_id"`
    Child_cat_id    string  `json:"child_cat_id"`
    Product_name    string  `json:"product_name"`
    Short_desc      string  `json:"short_desc"`
}

type ProductInput struct{
    User_id         int64           `json:"user_id"`
    Product_name    string          `json:"name"`
    Shop_id         int64           `json:"shop_id"`
    Child_cat_id    int64           `json:"category_id"`
    Short_desc      string          `json:"description"`
    Min_order       int64           `json:"min_order"`
    Price_currency  int64           `json:"price_currency"` // 1=IDR 2=USD
    Normal_price    int64           `json:"normal_price"`
    Status          int64           `json:"product_status"` // 
    Weight          ProductWeight   `json:"weight"`
    Insurance       int64           `json:"must_insurance"` // 0=no 1=yes
    Add_to_etalase  int64           `json:"add_to_etalase"` // 0=no 1=yes
    Etalase_id      int64           `json:"etalase_id"`
    Condition       int64           `json:"condition"`      // 1=baru 2=bekas
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

type ProductAlias struct{
    ProductId       int64           `bson:"product_id"`
    ProductKey      string          `bson:"product_key"`
    ShopId          int64           `bson:"shop_id"`
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

type Redis struct{
    Host        string
    Port        string
}
