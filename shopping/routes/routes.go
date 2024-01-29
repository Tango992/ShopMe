package routes

import (
	"net/http"
	"shopping/controller"
	_ "shopping/docs"
	"shopping/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

func Routes(e *echo.Echo, uc controller.UserController, pc controller.ProductController, tc controller.TransactionController) {
	users := e.Group("/users")
	{
		users.POST("/register", uc.Register)
		users.POST("/login", uc.Login)
		users.GET("/logout", uc.Logout)
	}
	
	products := e.Group("/products")
	{
		products.POST("", pc.Post)
		products.GET("", pc.GetAll)
		products.GET("/:productId", pc.GetById)
		products.PUT("/:productId", pc.UpdateById)
		products.DELETE("/:productId", pc.DeleteById)
	}

	transactions := e.Group("/transactions")
	transactions.Use(middlewares.RequireAuth)
	{
		transactions.POST("", tc.Create)
		transactions.GET("", tc.GetAll)
		transactions.GET("/:transactionId", tc.GetById)
		transactions.PUT("/:transactionId", tc.Update)
		transactions.DELETE("/:transactionId", tc.Delete)
	}

	e.GET("", func(c echo.Context) error {return c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")})
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}