package main

import (
	"AuthenticationService/internal/app"
	"log"
)

// @title Authentication Service
// @version 1.0

// @host localhost:8080
// @BasePath /

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	a.Run()
}
