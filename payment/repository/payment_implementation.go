package repository

import (
	"context"
	"payment/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentRepository struct {
	Collection *mongo.Collection
}

func NewPaymentRepository(collection *mongo.Collection) PaymentRepository {
	return PaymentRepository{
		Collection: collection,
	}
}

func (p PaymentRepository) Create(ctx context.Context, data *models.Payment) error {
	res, err := p.Collection.InsertOne(ctx, data)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	data.Id = res.InsertedID.(primitive.ObjectID)
	return nil
}