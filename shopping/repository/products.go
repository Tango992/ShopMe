package repository

import (
	"shopping/models"
	"shopping/utils"

	"go.mongodb.org/mongo-driver/bson"
)

type Products interface {
	Create(*models.Product) *utils.ErrResponse
	GetAll() ([]models.Product, *utils.ErrResponse)
	GetProductByName(string) (models.Product, *utils.ErrResponse)
	GetById(string) (models.Product, *utils.ErrResponse)
	Update(models.Product) *utils.ErrResponse
	Delete(string) (models.Product, *utils.ErrResponse)
	UpdateWithFilter(bson.M, bson.M) (int, *utils.ErrResponse) 
}