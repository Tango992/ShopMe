package repository

import (
	"context"
	"payment/models"

	"github.com/stretchr/testify/mock"
)

type MockPaymentRepository struct {
	Mock mock.Mock
}

func NewMockPaymentRepository() MockPaymentRepository {
	return MockPaymentRepository{}
}

func (m *MockPaymentRepository) Create(ctx context.Context, data *models.Payment) error {
	args := m.Mock.Called(ctx, data)
	return args.Error(0)
}