package main

import "eCommerce/storage/internal/application"

func main() {
	app := application.New()
	app.Build()
	app.InitTestProducts()
	app.Run()
}
