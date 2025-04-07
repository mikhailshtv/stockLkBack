package product

type Product struct {
	Id            int
	Code          int
	Name          string
	Quantity      string
	purchasePrice string
	SalePrice     string
}

func (product *Product) GetPurchasePrice() string {
	return product.purchasePrice
}

func (product *Product) SetPurchasePrice(purchasePrice string) {
	product.purchasePrice = purchasePrice
}
