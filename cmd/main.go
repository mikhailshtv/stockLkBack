package main

import (
	"context"
	"golang/stockLkBack/internal/app"
	"golang/stockLkBack/internal/config"
	"golang/stockLkBack/internal/service"
	"log"
)

func main() {
	service.RestoreData()

	newApp, err := app.NewApp(context.Background(), config.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	err = newApp.Start()
	if err != nil {
		log.Fatal(err)
	}
}
