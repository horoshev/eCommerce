package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrderUpdate struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty" extensions:"x-order=0"`
	Status    OrderStatus        `json:"status" bson:"status" extensions:"x-order=1"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp" extensions:"x-order=2"`
}
