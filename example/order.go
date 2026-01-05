package example

import (
	"context"
	"errors"
	"log"

	"hook-event-demo/event"
	"hook-event-demo/hook"
)

type Order struct {
	ID     string
	UserID string
	Amount int64
}

type Service struct {
	engine *hook.Engine
	bus    *event.Bus
}

func NewService(bus *event.Bus) *Service {
	engine := hook.NewEngine()

	engine.Register(hook.Before, &hook.Hook{
		Name:        "RiskCheck",
		Priority:    100,
		MustSucceed: true,
		Fn: func(ctx context.Context, data any) error {
			o := data.(*hook.OrderContext)
			if o.Amount > 10_000 {
				return errors.New("risk rejected")
			}
			return nil
		},
	})

	engine.Register(hook.After, &hook.Hook{
		Name:     "PublishOrderCreatedEvent",
		Priority: 10,
		Mode:     hook.Async,
		Fn: func(ctx context.Context, data any) error {
			o := data.(*hook.OrderContext)
			bus.Publish("order.created", event.OrderCreatedEvent{
				OrderID: o.OrderID,
				UserID:  o.UserID,
				Amount:  o.Amount,
			})
			return nil
		},
	})

	return &Service{engine: engine, bus: bus}
}

func (s *Service) CreateOrder(ctx context.Context, o *Order) error {
	ctxData := &hook.OrderContext{
		OrderID: o.ID,
		UserID:  o.UserID,
		Amount:  o.Amount,
	}

	if err := s.engine.Execute(ctx, hook.Before, ctxData); err != nil {
		return err
	}

	log.Printf("[CORE] order %s saved", o.ID)

	_ = s.engine.Execute(ctx, hook.After, ctxData)
	return nil
}
