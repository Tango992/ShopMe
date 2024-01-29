package main

import (
	"os"
	"shopping/config"
	"shopping/controller"
	"shopping/helpers"
	"shopping/repository"
	"shopping/routes"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title Graded Challenge 2
// @version 1.0
// @description Made for Graded Challenge 2 Hacktiv8 FTGO

// @contact.name Daniel Osvaldo Rahmanto
// @contact.email daniel.rahmanto@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	conn, paymentClient := config.InitGrpc()
	defer conn.Close()
	
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := config.ConnectDB().Database("graded1db")
	
	usersCollections := db.Collection("users")
	userRepository := repository.NewUserRepository(usersCollections)
	userController := controller.NewUserController(userRepository)
	
	productsCollection := db.Collection("products")
	productRepository := repository.NewProductRepository(productsCollection)
	productController := controller.NewProductController(productRepository)
	
	transactionCollection := db.Collection("transactions")
	transactionRepository := repository.NewTransactionRepository(transactionCollection)
	transactionController := controller.NewTransactionController(transactionRepository, productRepository, paymentClient)

	productController.StartCronJob()
	
	routes.Routes(e, userController, productController, transactionController)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}