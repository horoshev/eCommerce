package data

import (
	"context"
	"eCommerce/registry/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type RegistryRepository interface {
	FindOrderId(id primitive.ObjectID) (*models.Order, error)
	UpdateOrder(id primitive.ObjectID, updates []bson.E) (*models.Order, error)
	UpdateOrderStatus(id primitive.ObjectID, status models.OrderStatus) (*models.Order, error)
	Commit() error
}

type MongoRegistryRepository struct {
	orders *mongo.Collection
}

func NewMongoRegistryRepository(db *mongo.Database) *MongoRegistryRepository {
	r := new(MongoRegistryRepository)
	r.orders = db.Collection(`orders`)

	return r
}

func (m *MongoRegistryRepository) FindOrderId(id primitive.ObjectID) (*models.Order, error) {
	filter := bson.D{{"_id", id}}
	single := m.orders.FindOne(context.Background(), filter)
	if err := single.Err(); err != nil {
		return nil, err
	}

	order := new(models.Order)
	if err := single.Decode(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (m *MongoRegistryRepository) UpdateOrder(id primitive.ObjectID, updates []bson.E) (*models.Order, error) {
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", updates},
		{"$push", CreateStatusUpdateNote(updates)},
	}
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	single := m.orders.FindOneAndUpdate(context.Background(), filter, update, option)
	if err := single.Err(); err != nil {
		return nil, err
	}

	order := new(models.Order)
	if err := single.Decode(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (m *MongoRegistryRepository) UpdateOrderStatus(id primitive.ObjectID, status models.OrderStatus) (*models.Order, error) {
	return m.UpdateOrder(id, bson.D{{"status", status}})
}

func (m *MongoRegistryRepository) Commit() error {
	panic("implement me")
}

func CreateStatusUpdateNote(updates []bson.E) bson.M {
	for _, update := range updates {
		if update.Key != `status` {
			continue
		}

		status := update.Value.(models.OrderStatus)
		return bson.M{`updates`: models.OrderUpdate{
			Status:    status,
			Timestamp: time.Now().UTC(),
		}}
	}

	return bson.M{}

}
