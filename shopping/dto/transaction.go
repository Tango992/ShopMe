package dto

type Transaction struct {
	Product  string  `json:"product_name" validate:"required" extensions:"x-order=0"`
	Quantity uint    `json:"quantity" validate:"required,min=1" extensions:"x-order=1"`
}
