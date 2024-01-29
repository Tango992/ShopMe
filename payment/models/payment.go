package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payment struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StoreName   string             `bson:"store_name" json:"store_name"`
	CardHolder  string             `bson:"card_holder" json:"card_holder"`
	Amount      float32            `bson:"amount" json:"amount"`
	CompletedAt primitive.DateTime `bson:"completed" json:"completed"`
}
