package repository

import (
	"shopping/dto"
	"shopping/models"
	"shopping/utils"
)

type Users interface {
	Register(*models.User) *utils.ErrResponse
	FindUser(dto.LoginUser) (models.User, *utils.ErrResponse)
}