package repository

import (
	"context"
	"shopping/models"
	"shopping/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository struct {
	Collection *mongo.Collection
}

func NewTransactionRepository(collection *mongo.Collection) TransactionRepository {
	return TransactionRepository{
		Collection: collection,
	}
}

func (t TransactionRepository) Create(data *models.Transaction) *utils.ErrResponse {
	res, err1 := t.Collection.InsertOne(context.TODO(), data)
	if err1 != nil {
		return utils.ErrInternalServer.New(err1.Error())
	}
	
	data.Id = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (t TransactionRepository) GetAll(userEmail string) ([]models.Transaction, *utils.ErrResponse) {
	cursor, err := t.Collection.Find(context.TODO(), bson.M{"email": userEmail})
	if err != nil {
		return []models.Transaction{}, utils.ErrInternalServer.New(err.Error())
	}
	
	transactions := []models.Transaction{}
	if err := cursor.All(context.TODO(), &transactions); err != nil {
		return []models.Transaction{}, utils.ErrInternalServer.New(err.Error())
	}
	return transactions, nil
}

func (t TransactionRepository) GetById(userEmail, transId string) (models.Transaction, *utils.ErrResponse) {
	objectId, err := primitive.ObjectIDFromHex(transId)
	if err != nil {
		return models.Transaction{}, utils.ErrBadRequest.New(err.Error())
	}
	
	result := t.Collection.FindOne(context.TODO(), bson.M{"email": userEmail, "_id": objectId})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Transaction{}, utils.ErrNotFound.New(err.Error())
		}
		return models.Transaction{}, utils.ErrInternalServer.New(err.Error())
	}
	
	var transaction models.Transaction
	if err := result.Decode(&transaction); err != nil {
		return models.Transaction{}, utils.ErrInternalServer.New(err.Error())
	}
	return transaction, nil
}

func (t TransactionRepository) Update(data models.Transaction) *utils.ErrResponse {
	result := t.Collection.FindOneAndUpdate(context.TODO(), bson.M{"email": data.Email, "_id": data.Id}, bson.M{"$set": data})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrNotFound.New(err.Error())
		}
		return utils.ErrInternalServer.New(err.Error())
	}
	return nil
}

func (t TransactionRepository) Delete(userEmail, transId string) (models.Transaction, *utils.ErrResponse) {
	objectId, err := primitive.ObjectIDFromHex(transId)
	if err != nil {
		return models.Transaction{}, utils.ErrBadRequest.New(err.Error())
	}
	
	result := t.Collection.FindOneAndDelete(context.TODO(), bson.M{"email": userEmail, "_id": objectId})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Transaction{}, utils.ErrNotFound.New(err.Error())
		}
		return models.Transaction{}, utils.ErrInternalServer.New(err.Error())
	}
	
	var transaction models.Transaction
	if err := result.Decode(&transaction); err != nil {
		return models.Transaction{}, utils.ErrInternalServer.New(err.Error())
	}
	return transaction, nil
}