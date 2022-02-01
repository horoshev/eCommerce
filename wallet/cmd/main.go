package main

import "eCommerce/wallet/internal/application"

func main() {
	app := application.New()
	app.Build()
	app.Run()
}
