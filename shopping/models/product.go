package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id    primitive.ObjectID `bson:"_id,omitempty" json:"id" extensions:"x-order=0"`
	Name  string             `bson:"name" json:"name" extensions:"x-order=1"`
	Price float32            `bson:"price" json:"price" extensions:"x-order=2"`
	Stock uint               `bson:"stock" json:"stock" extensions:"x-order=3"`
}
