package messages

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}
