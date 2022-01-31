package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	Id     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Amount float64            `json:"amount"`
	Items  []OrderProduct     `json:"items" bson:"items"`
}

type OrderProduct struct {
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
}
