package repository

import (
	"context"
	"shopping/models"
	"shopping/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository struct {
	Collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) ProductRepository {
	return ProductRepository{
		Collection: collection,
	}
}

func (p ProductRepository) Create(data *models.Product) *utils.ErrResponse {
	productIndex := mongo.IndexModel{
		Keys: bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}
	
	if _, err := p.Collection.Indexes().CreateOne(context.TODO(), productIndex); err != nil {
		return utils.ErrInternalServer.New(err.Error())
	}
	
	res, err := p.Collection.InsertOne(context.TODO(), data)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return utils.ErrConflict.New(err.Error())
		}
		return utils.ErrInternalServer.New(err.Error())
	}
	
	data.Id = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (p ProductRepository) GetAll() ([]models.Product, *utils.ErrResponse) {
	cursor, err := p.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return []models.Product{}, utils.ErrInternalServer.New(err.Error())
	}
	
	products := []models.Product{}
	if err := cursor.All(context.TODO(), &products); err != nil {
		return []models.Product{}, utils.ErrInternalServer.New(err.Error())
	}
	return products, nil
}

func (p ProductRepository) GetById(productId string) (models.Product, *utils.ErrResponse) {
	objectId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return models.Product{}, utils.ErrBadRequest.New(err.Error())
	}

	result := p.Collection.FindOne(context.TODO(), bson.M{"_id": objectId})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Product{}, utils.ErrNotFound.New(err.Error())
		}
		return models.Product{}, utils.ErrInternalServer.New(err.Error())
	}
	
	var product models.Product
	if err := result.Decode(&product); err != nil {
		return models.Product{}, utils.ErrInternalServer.New(err.Error())
	}
	return product, nil
}

func (p ProductRepository) GetProductByName(productName string) (models.Product, *utils.ErrResponse) {
	result := p.Collection.FindOne(context.TODO(), bson.M{"name": productName})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Product{}, utils.ErrNotFound.New(err.Error())
		}
		return models.Product{}, utils.ErrInternalServer.New(err.Error())
	}
	
	var product models.Product
	if err := result.Decode(&product); err != nil {
		return models.Product{}, utils.ErrInternalServer.New(err.Error())
	}
	return product, nil
}

func (p ProductRepository) Update(data models.Product) *utils.ErrResponse {
	result := p.Collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": data.Id}, bson.M{"$set": data})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrNotFound.New(err.Error())
		}
		return utils.ErrInternalServer.New(err.Error())
	}
	return nil
}

func (p ProductRepository) Delete(productId string) (models.Product, *utils.ErrResponse) {
	objectId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return models.Product{}, utils.ErrBadRequest.New(err.Error())
	}

	result := p.Collection.FindOneAndDelete(context.TODO(), bson.M{"_id": objectId})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Product{}, utils.ErrNotFound.New(err.Error())
		}
		return models.Product{}, utils.ErrInternalServer.New(err.Error())
	}

	var product models.Product
	if err := result.Decode(&product); err != nil {
		return models.Product{}, utils.ErrInternalServer.New(err.Error())
	}
	return product, nil
}

func (p ProductRepository) UpdateWithFilter(filter, update bson.M) (int, *utils.ErrResponse) {
	result, err := p.Collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return 0, utils.ErrInternalServer.New(err.Error())
	}
	return int(result.ModifiedCount), nil
}