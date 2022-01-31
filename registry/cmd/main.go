package main

import (
	"eCommerce/registry/internal/application"
)

const (
	Version = `1.0.0`
)

// @title Registry
// @version 1.0
// @description This API is entrypoint for user requests.
// @description User has option to order some products and list created orders.
// @description Also, user is able to list his requests to this API.

// @host localhost:80
// @BasePath /
// @query.collection.format multi

// @x-extension-openapi {"example": "value on a json format"}
func main() {
	app := application.New()
	app.Build()
	app.Run()
}
