package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty" extensions:"x-order=0"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id" extensions:"x-order=1"`
	Status    OrderStatus        `json:"status" bson:"status" extensions:"x-order=2"`
	Amount    float64            `json:"amount" bson:"amount" extensions:"x-order=3"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp" extensions:"x-order=4"`
	Items     []OrderProduct     `json:"items" bson:"items" extensions:"x-order=5"`
	Updates   []OrderUpdate      `json:"-" bson:"updates" extensions:"x-order=6"`
}
