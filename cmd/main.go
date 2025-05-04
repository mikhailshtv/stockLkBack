package main

import (
	"golang/stockLkBack/internal/service"
	"time"
)

func main() {
	service.Interval()
	go service.LogAddedEntities()
	time.Sleep(time.Second * 20)
}
