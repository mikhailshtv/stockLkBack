package model

type Product struct {
	ID            int32  `json:"id"`
	Code          int32  `json:"code"`
	Quantity      int32  `json:"quantity"`
	Name          string `json:"name"`
	PurchasePrice int32  `json:"purchasePrice"`
	SalePrice     int32  `json:"salePrice"`
}

type ProductRequestBody struct {
	Code          int32  `json:"code"`
	Quantity      int32  `json:"quantity"`
	Name          string `json:"name"`
	PurchasePrice int32  `json:"purchasePrice"`
	SalePrice     int32  `json:"salePrice"`
}
