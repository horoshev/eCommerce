package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Wallet struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	UserId       primitive.ObjectID `bson:"user_id"`
	UserName     string             `bson:"user_name"`
	Balance      float64            `bson:"balance"`
	Transactions []Transaction      `bson:"transactions"`
}
