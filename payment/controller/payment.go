package controller

import (
	"context"
	"payment/models"
	"payment/pb"
	"payment/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentContoller struct {
	pb.UnimplementedPaymentServer
	Collection repository.Payment
}

func NewPaymentController(collection repository.Payment) PaymentContoller {
	return PaymentContoller{
		Collection: collection,
	}
}

func (p PaymentContoller) Create(ctx context.Context, data *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	currentTime := time.Now()
	paymentData := models.Payment{
		StoreName: data.StoreName,
		CardHolder: data.CardHolder,
		Amount: data.Amount,
		CompletedAt: primitive.NewDateTimeFromTime(currentTime),
	}
	
	if err := p.Collection.Create(ctx, &paymentData); err != nil {
		return nil, err
	}

	paymentResponse := &pb.PaymentResponse{
		Id: paymentData.Id.Hex(),
		StoreName: paymentData.StoreName,
		CardHolder: paymentData.CardHolder,
		Amount: paymentData.Amount,
		CompletedAt: currentTime.Format("2006-01-02 15:04:05"),
	}
	return paymentResponse, nil
}