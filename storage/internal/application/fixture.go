package application

import "eCommerce/storage/internal/models"

const (
	DefaultOffset   = 'A'
	DefaultQuantity = 10
)

func AsInterfaces(products []models.Product) []interface{} {
	out := make([]interface{}, len(products))

	for i := range products {
		out[i] = products[i]
	}

	return out
}

func GenerateProducts(n int) []models.Product {
	products := make([]models.Product, n)

	for i := 0; i < n; i += 1 {
		products[i] = NextProduct(i)
	}

	return products
}

func NextProduct(n int) models.Product {
	product := models.Product{
		Name:     GenerateProductName(n),
		Cost:     float64(n + 1),
		Quantity: DefaultQuantity,
	}

	return product
}

func GenerateProductName(n int) string {
	return string(DefaultOffset + n)
}
