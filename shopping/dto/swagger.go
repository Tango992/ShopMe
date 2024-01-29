package dto

import "shopping/models"

type GeneralResponse struct {
	Message string `json:"message" extensions:"x-order=0"`
	Data    string `json:"data" extensions:"x-order=1"`
}

type UserDataWithoutPassword struct {
	Id    string `bjson:"id" extensions:"x-order=0"`
	Name  string `json:"name" extensions:"x-order=1"`
	Email string `json:"email" extensions:"x-order=2"`
}

type RegisterResponse struct {
	Message string                  `json:"message" extensions:"x-order=0"`
	Data    UserDataWithoutPassword `json:"data" extensions:"x-order=1"`
}

type TransactionResponse struct {
	Message string             `json:"message" extensions:"x-order=0"`
	Data    models.Transaction `json:"data" extensions:"x-order=1"`
}

type TransactionsResponse struct {
	Message string               `json:"message" extensions:"x-order=0"`
	Data    []models.Transaction `json:"data" extensions:"x-order=1"`
}

type ProductResponse struct {
	Message string         `json:"message" extensions:"x-order=0"`
	Data    models.Product `json:"data" extensions:"x-order=1"`
}

type ProductsResponse struct {
	Message string           `json:"message" extensions:"x-order=0"`
	Data    []models.Product `json:"data" extensions:"x-order=1"`
}
