package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id" extensions:"x-order=0"`
	PaymentId primitive.ObjectID `bson:"payment_id" json:"payment_id" extensions:"x-order=1"`
	Email     string             `bson:"email" json:"email" extensions:"x-order=2"`
	Product   string             `bson:"product" json:"product" extensions:"x-order=3"`
	Quantity  uint               `bson:"quantity" json:"quantity" extensions:"x-order=4"`
	Total     float32            `bson:"total" json:"total" extensions:"x-order=5"`
}
