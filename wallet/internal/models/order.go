package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	Id     primitive.ObjectID `json:"id"`
	UserId primitive.ObjectID `json:"user_id"`
	Amount float64            `json:"amount"`
}
