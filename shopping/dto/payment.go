package dto

type Payment struct {
	StoreName  string  `json:"store_name" validate:"required"`
	CardHolder string  `json:"card_holder"  validate:"required"`
	Amount     float32 `json:"amount" validate:"required"`
}
