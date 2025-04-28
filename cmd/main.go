package main

import (
	"golang/stockLkBack/internal/repository"
	"golang/stockLkBack/internal/service"
	"sync"
	"time"
)

func main() {
	go inerval()
	go repository.LogAddedEntities()
	time.Sleep(time.Second * 20)
}

func inerval() {
	for range time.Tick(time.Second * 1) {
		channel := make(chan any)
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			channel <- service.NewEntity()
			wg.Done()
		}()
		go func() {
			repository.CheckAndSaveEntity(<-channel)
			wg.Done()
		}()
		wg.Wait()
	}
}
