package model

type Product struct {
	Id            int
	Code          int
	Quantity      int
	Name          string
	purchasePrice int
	SalePrice     int
}

func (product *Product) PurchasePrice() int {
	return product.purchasePrice
}

func (product *Product) SetPurchasePrice(purchasePrice int) {
	product.purchasePrice = purchasePrice
}
