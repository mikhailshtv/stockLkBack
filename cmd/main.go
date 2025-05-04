package main

import (
	"context"
	"golang/stockLkBack/internal/service"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup
	service.Interval(ctx, &wg)
	service.LogAddedEntities(ctx, &wg)
	wg.Wait()
}
