package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Cost     float64            `bson:"cost"`
	Quantity int64              `bson:"quantity"`
}
