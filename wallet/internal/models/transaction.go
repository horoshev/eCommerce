package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	TransactionActive    TransactionStatus = `TRANSACTION_ACTIVE`
	TransactionCancelled TransactionStatus = `TRANSACTION_CANCELLED`
)

type TransactionStatus string

type Transaction struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	OrderId   primitive.ObjectID `bson:"order_id"`
	Status    TransactionStatus  `bson:"status"`
	Amount    float64            `bson:"amount"`
	Balance   float64            `bson:"balance"`
	Note      string             `bson:"note"`
	Timestamp time.Time          `bson:"timestamp"`
}
