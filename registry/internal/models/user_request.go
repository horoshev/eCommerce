package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserRequest struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty" extensions:"x-order=0"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id" extensions:"x-order=1"`
	Type      string             `json:"type" bson:"type" extensions:"x-order=2"`
	Path      string             `json:"path" bson:"path" extensions:"x-order=3"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp" extensions:"x-order=4"`
}
