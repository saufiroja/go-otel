package main

import "github.com/saufiroja/go-otel/auth-service/internal/app"

func main() {
	apps := app.NewApp()
	apps.Start()
}
