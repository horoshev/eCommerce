package requests

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderRequest struct {
	OrderId  primitive.ObjectID `json:"order_id"`
	Products []OrderProduct     `json:"products"`
}

type OrderProduct struct {
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
}
