package controller

import (
	"fmt"
	"net/http"
	"shopping/dto"
	"shopping/models"
	"shopping/repository"
	"shopping/utils"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductController struct {
	Cron *cron.Cron
	Repository repository.Products
}

func NewProductController(r repository.Products) ProductController {
	return ProductController{
		Cron: cron.New(),
		Repository: r,
	}
}

// Create        godoc
// @Summary      Add new product into database
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request body dto.Product  true  "Product data"
// @Success      201  {object}  dto.ProductResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /products [post]
func (p ProductController) Post(c echo.Context) error {
	var productDataTmp dto.Product
	if err := c.Bind(&productDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}
	
	if err := c.Validate(&productDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	productData := models.Product{
		Name: productDataTmp.Name,
		Price: productDataTmp.Price,
		Stock: productDataTmp.Stock,
	}
	
	if err := p.Repository.Create(&productData); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}
	
	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Product added",
		Data: productData,
	})
}

// Get all       godoc
// @Summary      Get all products from database
// @Tags         products
// @Produce      json
// @Success      200  {object}  dto.ProductsResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /products [get]
func (p ProductController) GetAll(c echo.Context) error {
	products, err := p.Repository.GetAll()
	if err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "Get all products",
		Data: products,
	})
}

// Get by ID     godoc
// @Summary      Get specific product from database
// @Tags         products
// @Produce      json
// @Param        productId   path   string   true   "Product ID"
// @Success      200  {object}  dto.ProductResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      404  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /products/{productId} [get]
func (p ProductController) GetById(c echo.Context) error {
	productId := c.Param("productId")
	product, err := p.Repository.GetById(productId)
	if err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "Get product by id",
		Data: product,
	})
}

// Update product     godoc
// @Summary      Update specific product from database
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        productId   path   string   true   "Product ID"
// @Param        request body dto.Product  true  "Updated product data"
// @Success      200  {object}  dto.ProductResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      404  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /products/{productId} [put]
func (p ProductController) UpdateById(c echo.Context) error {
	productId := c.Param("productId")
	objectId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}
	
	var productDataTmp dto.Product
	if err := c.Bind(&productDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}
	
	if err := c.Validate(&productDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	productData := models.Product{
		Id: objectId,
		Name: productDataTmp.Name,
		Price: productDataTmp.Price,
		Stock: productDataTmp.Stock,
	}
	
	if err := p.Repository.Update(productData); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "Product updated",
		Data: productData,
	})
}

// Delete     godoc
// @Summary      Delete specific product from database
// @Tags         products
// @Produce      json
// @Param        productId   path   string   true   "Product ID"
// @Success      200  {object}  dto.ProductResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      404  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /products/{productId} [delete]
func (p ProductController) DeleteById(c echo.Context) error {
	productId := c.Param("productId")
	product, err := p.Repository.Delete(productId)
	if err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "Product deleted",
		Data: product,
	})
}

func (p ProductController) RefreshStock() {
	p.Cron.AddFunc("0 0 * * *", func() {
		filter := bson.M{}
		update := bson.M{"$set": bson.M{"stock": 10000}}

		updatedCount, err := p.Repository.UpdateWithFilter(filter, update)
		if err != nil {
			fmt.Printf("Crob job failed: %v", err)
		}
		fmt.Printf("Cron job updated %v entries\n", updatedCount)
	})
}

func (p ProductController) StartCronJob() {
	p.RefreshStock()
	p.Cron.Start()
}