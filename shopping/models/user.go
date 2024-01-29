package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id" extensions:"x-order=0"`
	Name     string             `bson:"name" json:"name" extensions:"x-order=1"`
	Email    string             `bson:"email" json:"email" extensions:"x-order=2"`
	Password string             `bson:"password" json:"password,omitempty" extensions:"x-order=3"`
}
