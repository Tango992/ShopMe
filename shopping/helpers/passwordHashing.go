package helpers

import (
	"shopping/dto"
	"shopping/models"
	"shopping/utils"

	"golang.org/x/crypto/bcrypt"
)

func CreateHash(data *dto.RegisterUser) *utils.ErrResponse {
	hashed, err:= bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrInternalServer.New(err.Error())
	}
	data.Password = string(hashed)
	return nil
}

func CheckPassword(dbData models.User, data dto.LoginUser) *utils.ErrResponse {
	if err := bcrypt.CompareHashAndPassword([]byte(dbData.Password), []byte(data.Password)); err != nil {
		return utils.ErrUnauthorized.New("Invalid email/password")
	}
	return nil
}