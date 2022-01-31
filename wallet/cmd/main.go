package main

import "eCommerce/wallet/internal/application"

const ServiceName = application.ServiceName

//const Env = application.Production
const Env = application.Development

func main() {
	app := application.New(Env)
	app.Build()
	app.Run()
}
