package models

type OrderProduct struct {
	Name     string `json:"name" bson:"name" extensions:"x-order=0"`
	Quantity int64  `json:"quantity" bson:"quantity" extensions:"x-order=1"`
}
