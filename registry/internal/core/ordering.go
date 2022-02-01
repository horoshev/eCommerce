package core

import (
	"context"
	"eCommerce/registry/internal/api/requests"
	"eCommerce/registry/internal/models"
	"errors"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type PurchaseController interface {
	Order(userId primitive.ObjectID, r *requests.OrderRequest) (*models.Order, error)
	ListOrders(r *requests.PageRequest) ([]models.Order, error)
	ListUserOrders(userId primitive.ObjectID, r *requests.PageRequest) ([]models.Order, error)
}

type Purchaser struct {
	log         *zap.SugaredLogger
	Orders      *mongo.Collection
	Producer    *kafka.Writer
	Coordinator *OrderCoordinator
}

func NewPurchaser(log *zap.SugaredLogger, db *mongo.Database, producer *kafka.Writer, c *OrderCoordinator) *Purchaser {
	p := new(Purchaser)
	p.log = log
	p.Orders = db.Collection("orders")
	p.Producer = producer
	p.Coordinator = c

	return p
}

// Order creating an order and publish event
func (p *Purchaser) Order(userId primitive.ObjectID, r *requests.OrderRequest) (*models.Order, error) {
	if r.Items == nil || len(r.Items) == 0 {
		return nil, errors.New(`no items in order`)
	}

	order := new(models.Order)
	order.UserId = userId
	order.Status = models.OrderPending
	order.Timestamp = time.Now().UTC()
	order.Items = r.Items
	order.Updates = []models.OrderUpdate{
		{
			Status:    models.OrderPending,
			Timestamp: time.Now().UTC(),
		},
	}

	one, err := p.Orders.InsertOne(context.Background(), order)
	if err != nil {
		return nil, err
	}

	order.Id = one.InsertedID.(primitive.ObjectID)
	p.Coordinator.NewOrder(order)

	return order, nil
}

func (p *Purchaser) ListOrders(r *requests.PageRequest) ([]models.Order, error) {
	opt := options.Find()
	opt.SetSort(bson.D{{"_id", -1}})
	opt.SetSkip(int64(r.Page * r.Size))
	opt.SetLimit(int64(r.Size))

	records, err := p.Orders.Find(context.Background(), bson.D{}, opt)
	if err != nil {
		return nil, err
	}

	list := make([]models.Order, 0, r.Size)
	for records.Next(context.Background()) {
		var record models.Order
		if err = records.Decode(&record); err != nil {
			return nil, err
		}
		list = append(list, record)
	}

	return list, nil
}

func (p *Purchaser) ListUserOrders(userId primitive.ObjectID, r *requests.PageRequest) ([]models.Order, error) {
	opt := options.Find()
	opt.SetSort(bson.D{{"_id", -1}})
	opt.SetSkip(int64(r.Page * r.Size))
	opt.SetLimit(int64(r.Size))
	filter := bson.D{{"user_id", userId}}
	records, err := p.Orders.Find(context.Background(), filter, opt)
	if err != nil {
		return nil, err
	}

	list := make([]models.Order, 0, r.Size)
	for records.Next(context.Background()) {
		var record models.Order
		if err = records.Decode(&record); err != nil {
			return nil, err
		}
		list = append(list, record)
	}

	return list, nil
}
