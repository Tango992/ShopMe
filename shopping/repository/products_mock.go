package repository

import (
	"shopping/models"
	"shopping/utils"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockProductRepository struct {
	Mock mock.Mock
}

func NewMockProductRepository() MockProductRepository {
	return MockProductRepository{}
}

func (m *MockProductRepository) Create(data *models.Product) *utils.ErrResponse {
	args := m.Mock.Called(data)
	args.Get(0).(*models.Product).Id = primitive.NewObjectID()
	return args.Get(1).(*utils.ErrResponse)
}

func (m *MockProductRepository) GetAll() ([]models.Product, *utils.ErrResponse) {
	args := m.Mock.Called()
	return args.Get(0).([]models.Product), args.Get(1).(*utils.ErrResponse)
}

func (m *MockProductRepository) GetProductByName(product string) (models.Product, *utils.ErrResponse) {	
	args := m.Mock.Called(product)
	return args.Get(0).(models.Product), args.Get(1).(*utils.ErrResponse)
}

func (m *MockProductRepository) GetById(productId string) (models.Product, *utils.ErrResponse) {
	args := m.Mock.Called(productId)
	return args.Get(0).(models.Product), args.Get(1).(*utils.ErrResponse)
}

func (m *MockProductRepository) Update(data models.Product) *utils.ErrResponse {
	args := m.Mock.Called(data)
	return args.Get(0).(*utils.ErrResponse)
}

func (m *MockProductRepository) Delete(productId string) (models.Product, *utils.ErrResponse) {
	args := m.Mock.Called(productId)
	return args.Get(0).(models.Product), args.Get(1).(*utils.ErrResponse)
}

func (m *MockProductRepository) UpdateWithFilter(filter, update bson.M) (int, *utils.ErrResponse) {
	args := m.Mock.Called(filter, update)
	return args.Int(0), args.Get(1).(*utils.ErrResponse)
}