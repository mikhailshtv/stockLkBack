package model

type Product struct {
	ID            int32  `json:"id" db:"id"`
	Code          int32  `json:"code" db:"code"`
	Quantity      int32  `json:"quantity" db:"quantity"`
	Name          string `json:"name" db:"name"`
	PurchasePrice int32  `json:"purchasePrice,omitempty" db:"purchase_price"`
	SellPrice     int32  `json:"sellPrice" db:"sell_price"`
	Version       int    `json:"-" db:"version"`
}

type ProductRequestBody struct {
	Code          int32  `json:"code" db:"code"`
	Quantity      int32  `json:"quantity" db:"quantity"`
	Name          string `json:"name" db:"name"`
	PurchasePrice int32  `json:"purchasePrice" db:"purchase_price"`
	SellPrice     int32  `json:"sellPrice" db:"sell_price"`
}

// ProductQueryParams параметры запроса для списка продуктов
// @Description Параметры запроса для фильтрации, сортировки и пагинации списка продуктов.
type ProductQueryParams struct {
	Code          *int   `form:"code" json:"code,omitempty" example:"123"`
	Quantity      *int   `form:"quantity" json:"quantity,omitempty" example:"10"`
	Name          string `form:"name" json:"name,omitempty" example:"Молоко"`
	PurchasePrice *int   `form:"purchase_price" json:"purchasePrice,omitempty" example:"100"`
	SellPrice     *int   `form:"sell_price" json:"sellPrice,omitempty" example:"150"`
	SortField     string `form:"sort_field" json:"sortField,omitempty" example:"name"`
	SortOrder     string `form:"sort_order" json:"sortOrder,omitempty" example:"ASC"`
	Page          int    `form:"page" json:"page,omitempty" example:"1"`
	PageSize      int    `form:"page_size" json:"pageSize,omitempty" example:"10"`
}

// ProductListResponse ответ со списком продуктов
// @Description Ответ со списком продуктов и метаданными пагинации.
type ProductListResponse struct {
	Data     []Product `json:"data"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
	Total    int       `json:"total,omitempty"`
}
