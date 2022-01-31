package models

type OrderStatus string

const (
	OrderError                    OrderStatus = `ORDER_ERROR`
	OrderPending                  OrderStatus = `ORDER_PENDING`
	OrderReservationPending       OrderStatus = `ORDER_RESERVATION_PENDING`
	OrderReserved                 OrderStatus = `ORDER_RESERVED`
	OrderPaymentPending           OrderStatus = `ORDER_PAYMENT_PENDING`
	OrderPaid                     OrderStatus = `ORDER_PAID`
	OrderCancelPending            OrderStatus = `ORDER_CANCEL_PENDING`
	OrderPaymentCancelPending     OrderStatus = `ORDER_PAYMENT_CANCEL_PENDING`
	OrderPaymentCanceled          OrderStatus = `ORDER_PAYMENT_CANCELED`
	OrderReservationCancelPending OrderStatus = `ORDER_RESERVATION_CANCEL_PENDING`
	OrderReservationCanceled      OrderStatus = `ORDER_RESERVATION_CANCELED`
	OrderCanceled                 OrderStatus = `ORDER_CANCELED`
	OrderCancellationError        OrderStatus = `ORDER_CANCELLATION_ERROR`
)
