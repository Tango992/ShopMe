package controller

import (
	"context"
	"fmt"
	"payment/models"
	"payment/pb"
	"payment/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	mockRepository = repository.NewMockPaymentRepository()
	paymentController = NewPaymentController(&mockRepository)
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestCreate(t *testing.T) {
	dummyObjectId := primitive.NewObjectID()
	
	repositoryRequestData := &models.Payment{
		StoreName: "IndoApril Cabang 12345",
		CardHolder: "Windah Basudara",
		Amount: 100000,
		CompletedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	grpcRequestData := &pb.PaymentRequest{
		StoreName: "IndoApril Cabang 12345",
		CardHolder: "Windah Basudara",
		Amount: 100000,
	}

	mockRepository.Mock.On("Create", context.TODO(), repositoryRequestData).Return(nil).Run(func(args mock.Arguments) {
		paymentData := args.Get(1).(*models.Payment)
		paymentData.Id = dummyObjectId
	})

	res, err := paymentController.Create(context.TODO(), grpcRequestData)
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.Equal(t, dummyObjectId.Hex(), res.GetId())
	
	fmt.Println(res)
}