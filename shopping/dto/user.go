package dto

type RegisterUser struct {
	Name     string `json:"name" validate:"required" extensions:"x-order=0"`
	Email    string `json:"email" validate:"required" extensions:"x-order=1"`
	Password string `json:"password" validate:"required" extensions:"x-order=2"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required" extensions:"x-order=0"`
	Password string `json:"password" validate:"required" extensions:"x-order=1"`
}
