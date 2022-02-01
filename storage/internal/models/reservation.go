package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Success  ReservationStatus = `success`
	Canceled ReservationStatus = `canceled`
)

type ReservationStatus string

type OrderReservation struct {
	Id       primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	OrderId  primitive.ObjectID   `json:"order_id" bson:"order_id"`
	Products []ProductReservation `json:"products" bson:"products"`
	Status   ReservationStatus    `json:"status" bson:"status"`
}

type ProductReservation struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductName string             `json:"product_name" bson:"product_name"`
	Quantity    int64              `json:"quantity" bson:"quantity"`
}

func OrderToReservation(order *Order) *OrderReservation {
	r := new(OrderReservation)
	r.Status = Success
	r.OrderId = order.Id
	r.Products = make([]ProductReservation, len(order.Items))

	for i, x := range order.Items {
		r.Products[i] = ProductReservation{
			ProductName: x.Name,
			Quantity:    x.Quantity,
		}
	}

	return r
}
