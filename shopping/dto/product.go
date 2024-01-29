package dto

type Product struct {
	Name  string  `json:"name" validate:"required" extensions:"x-order=0"`
	Price float32 `json:"price" validate:"required" extensions:"x-order=1"`
	Stock uint    `json:"stock" validate:"required" extensions:"x-order=2"`
}
