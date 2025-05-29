package model

type Product struct {
	Id            int    `json:"id"`
	Code          int    `json:"code"`
	Quantity      int    `json:"quantity"`
	Name          string `json:"name"`
	PurchasePrice int    `json:"purchasePrice"`
	SalePrice     int    `json:"salePrice"`
}

// func (product *Product) PurchasePrice() int {
// 	return product.purchasePrice
// }

// func (product *Product) SetPurchasePrice(purchasePrice int) {
// 	product.purchasePrice = purchasePrice
// }
