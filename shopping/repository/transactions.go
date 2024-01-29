package repository

import (
	"shopping/models"
	"shopping/utils"
)

type Transactions interface{
	Create(*models.Transaction) *utils.ErrResponse
	GetAll(string) ([]models.Transaction, *utils.ErrResponse)
	GetById(string, string) (models.Transaction, *utils.ErrResponse)
	Update(models.Transaction) *utils.ErrResponse
	Delete(string, string) (models.Transaction, *utils.ErrResponse)
}