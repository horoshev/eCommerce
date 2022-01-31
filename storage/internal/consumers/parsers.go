package consumers

import (
	"eCommerce/storage/internal/models"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

func ParseOrder(message kafka.Message) (*models.Order, error) {
	order := new(models.Order)
	err := json.Unmarshal(message.Value, order)

	if err != nil {
		return nil, err
	}

	return order, nil
}
