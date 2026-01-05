package main

import (
	"context"
	"log"
	"time"

	"hook-event-demo/event"
	"hook-event-demo/example"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	bus := event.NewBus()
	event.RegisterConsumers(bus)

	service := example.NewService(bus)

	log.Println("=== create normal order ===")
	_ = service.CreateOrder(context.Background(), &example.Order{
		ID:     "order_1001",
		UserID: "user_1",
		Amount: 500,
	})

	log.Println("=== create risk order ===")
	err := service.CreateOrder(context.Background(), &example.Order{
		ID:     "order_1002",
		UserID: "user_2",
		Amount: 20_000,
	})
	if err != nil {
		log.Println("order failed:", err)
	}

	time.Sleep(time.Second)
}
