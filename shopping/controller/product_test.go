package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"shopping/dto"
	"shopping/helpers"
	"shopping/models"
	"shopping/repository"
	"shopping/utils"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	e *echo.Echo = echo.New()
	mockRepository = repository.NewMockProductRepository()
	productController = NewProductController(&mockRepository)
)

func TestMain(m *testing.M) {
	e = echo.New()
	e.Validator = &helpers.CustomValidator{NewValidator: validator.New()}

	m.Run()
}

func TestCreate(t *testing.T) {
	requestBody := dto.Product{
		Name: "Teh Kotak",
		Price: 4000,
		Stock: 1000,
	}
	requestBodyBytes, _ := json.Marshal(requestBody)
	
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/products", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	
	productData := &models.Product{
		Name: "Teh Kotak",
		Price: 4000,
		Stock: 1000,
	}

	mockRepository.Mock.On("Create", mock.Anything).Return(productData, (*utils.ErrResponse)(nil)).Run(func(args mock.Arguments) {
		product := args.Get(0).(*models.Product)
		product.Id = primitive.NewObjectID()
	})

	productController.Post(c)
	response := w.Result()
	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(string(body))

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotEmpty(t, body)
	assert.Equal(t, "Product added", responseBody["message"].(string))
}

func TestGetAll(t *testing.T) {
	productDatas := []models.Product{{
		Id: primitive.NewObjectID(),
		Name: "Teh Kotak",
		Price: 4000,
		Stock: 1000,
		},{
		Id: primitive.NewObjectID(),
		Name: "Teh Botol",
		Price: 3500,
		Stock: 1000,
	},}
	
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/products", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	
	mockRepository.Mock.On("GetAll").Return(productDatas, (*utils.ErrResponse)(nil))

	productController.GetAll(c)
	response := w.Result()
	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(string(body))

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotEmpty(t, body)
	assert.Equal(t, "Get all products", responseBody["message"].(string))
}

func TestGetById(t *testing.T) {
	mockProductId := "655a1d0c98d1183ed4045011"
	objectId, _ := primitive.ObjectIDFromHex(mockProductId)

	productData := models.Product{
		Id: objectId,
		Name: "Teh Kotak",
		Price: 4000,
		Stock: 1000,
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/products/:productId", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	
	c.SetParamNames("productId")
	c.SetParamValues(mockProductId)

	mockRepository.Mock.On("GetById", mock.Anything).Return(productData, (*utils.ErrResponse)(nil))

	productController.GetById(c)
	response := w.Result()
	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(string(body))

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotEmpty(t, body)
	assert.Equal(t, "Get product by id", responseBody["message"].(string))
}

func TestUpdate(t *testing.T) {
	mockProductId := "655a1d0c98d1183ed4045011"
	objectId, _ := primitive.ObjectIDFromHex(mockProductId)
	
	productData := models.Product{
		Id: objectId,
		Name: "Teh Kotak",
		Price: 4000,
		Stock: 1000,
	}
	requestBodyBytes, _ := json.Marshal(productData)
	
	req := httptest.NewRequest(http.MethodPut, "http://localhost:8080/products/:productId", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c := e.NewContext(req, w)

	c.SetParamNames("productId")
	c.SetParamValues(mockProductId)

	mockRepository.Mock.On("Update", mock.Anything).Return((*utils.ErrResponse)(nil))

	productController.UpdateById(c)
	response := w.Result()
	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(string(body))

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotEmpty(t, body)
	assert.Equal(t, "Product updated", responseBody["message"].(string))
}

func TestDelete(t *testing.T) {
	mockProductId := "655a1d0c98d1183ed4045011"
	objectId, _ := primitive.ObjectIDFromHex(mockProductId)
	
	productData := models.Product{
		Id: objectId,
		Name: "Teh Kotak",
		Price: 4000,
		Stock: 1000,
	}
	
	req := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/products/:productId", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c := e.NewContext(req, w)

	c.SetParamNames("productId")
	c.SetParamValues(mockProductId)

	mockRepository.Mock.On("Delete", mock.Anything).Return(productData, (*utils.ErrResponse)(nil))

	productController.DeleteById(c)
	response := w.Result()
	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(string(body))

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotEmpty(t, body)
	assert.Equal(t, "Product deleted", responseBody["message"].(string))
}