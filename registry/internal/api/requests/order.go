package requests

import (
	"eCommerce/registry/internal/models"
	"time"
)

type OrderRequest struct {
	Items []models.OrderProduct `json:"items"`
}

type OrderPageRequest struct {
	PageRequest
	Methods []string  `json:"methods"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
}
