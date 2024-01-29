package controller

import (
	"net/http"
	"shopping/dto"
	"shopping/helpers"
	"shopping/models"
	"shopping/repository"
	"shopping/utils"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	Repository repository.Users
}

func NewUserController(r repository.Users) UserController {
	return UserController{
		Repository: r,
	}
}

// Register      godoc
// @Summary      Register new user into database
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterUser  true  "Register data"
// @Success      201  {object}  dto.RegisterResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      409  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /users/register [post]
func (u UserController) Register(c echo.Context) error {
	var registerDataTmp dto.RegisterUser
	if err := c.Bind(&registerDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}
	
	if err := c.Validate(&registerDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	if err := helpers.CreateHash(&registerDataTmp); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}

	registerData := models.User{
		Name: registerDataTmp.Name,
		Email: registerDataTmp.Email,
		Password: registerDataTmp.Password,
	}
	
	if err := u.Repository.Register(&registerData); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}
	
	registerData.Password = ""
	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Registered",
		Data: registerData,
	})
}

// Login         godoc
// @Summary      Log in with existing account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginUser  true  "Login data"
// @Success      200  {object}  dto.GeneralResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /users/login [post]
func (u UserController) Login(c echo.Context) error {
	var loginData dto.LoginUser
	if err := c.Bind(&loginData); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}
	
	if err := c.Validate(&loginData); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	dbData, err := u.Repository.FindUser(loginData)
	if err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}
	
	if err := helpers.CheckPassword(dbData, loginData); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}

	if err := helpers.SignNewJWT(c, dbData); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}
	
	return c.JSON(http.StatusOK, dto.Response{
		Message: "Logged in",
		Data: "Authorization token is stored using cookie",
	})
}

// Logout        godoc
// @Summary      Clears authorization cookie
// @Tags         users
// @Produce      json
// @Success      200  {object}  dto.GeneralResponse
// @Router       /users/logout [get]
func (u UserController) Logout(c echo.Context) error {
	cookie := new(http.Cookie)

	cookie.Name = "Authorization"
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Value = ""
	cookie.SameSite = http.SameSiteLaxMode
	cookie.MaxAge = -1
	c.SetCookie(cookie)
	
	return c.JSON(http.StatusOK, dto.Response{
		Message: "Logged out",
		Data: "Authorization cookie has been deleted",
	})
}