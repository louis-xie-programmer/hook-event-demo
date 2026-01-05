package event

import "log"

func RegisterConsumers(bus *Bus) {
	bus.Subscribe("order.created", func(e any) {
		o := e.(OrderCreatedEvent)
		log.Printf("[SMS] order=%s", o.OrderID)
	})

	bus.Subscribe("order.created", func(e any) {
		o := e.(OrderCreatedEvent)
		log.Printf("[POINT] user=%s", o.UserID)
	})

	bus.Subscribe("order.created", func(e any) {
		o := e.(OrderCreatedEvent)
		log.Printf("[BI] order=%s amount=%d", o.OrderID, o.Amount)
	})
}
