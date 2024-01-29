package main

import (
	"payment/config"
	"payment/controller"
	"payment/repository"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db := config.ConnectDB().Database("graded1db")
	paymentCollection := db.Collection("payments")
	
	paymentRepository := repository.NewPaymentRepository(paymentCollection)
	paymentController := controller.NewPaymentController(paymentRepository)

	config.ListenAndServeGrpc(paymentController)
}