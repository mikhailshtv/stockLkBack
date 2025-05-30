package model

type Product struct {
	Id            int    `json:"id"`
	Code          int    `json:"code"`
	Quantity      int    `json:"quantity"`
	Name          string `json:"name"`
	PurchasePrice int    `json:"purchasePrice"`
	SalePrice     int    `json:"salePrice"`
}

type ProductRequestBody struct {
	Code          int    `json:"code"`
	Quantity      int    `json:"quantity"`
	Name          string `json:"name"`
	PurchasePrice int    `json:"purchasePrice"`
	SalePrice     int    `json:"salePrice"`
}
