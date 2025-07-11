package model

type Product struct {
	ID            int32  `json:"id" db:"id"`
	Code          int32  `json:"code" db:"code"`
	Quantity      int32  `json:"quantity" db:"quantity"`
	Name          string `json:"name" db:"name"`
	PurchasePrice int32  `json:"purchasePrice,omitempty" db:"purchase_price"`
	SellPrice     int32  `json:"sellPrice" db:"sell_price"`
}

type ProductRequestBody struct {
	Code          int32  `json:"code" db:"code"`
	Quantity      int32  `json:"quantity" db:"quantity"`
	Name          string `json:"name" db:"name"`
	PurchasePrice int32  `json:"purchasePrice" db:"purchase_price"`
	SellPrice     int32  `json:"sellPrice" db:"sell_price"`
}
