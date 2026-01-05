package hook

import (
	"context"
	"fmt"
	"sort"
)

type Type int

const (
	Before Type = iota
	After
)

type Engine struct {
	hooks map[Type][]*Hook
}

func NewEngine() *Engine {
	return &Engine{hooks: make(map[Type][]*Hook)}
}

func (e *Engine) Register(t Type, h *Hook) {
	e.hooks[t] = append(e.hooks[t], h)
	sort.Slice(e.hooks[t], func(i, j int) bool {
		return e.hooks[t][i].Priority > e.hooks[t][j].Priority
	})
}

func (e *Engine) Execute(ctx context.Context, t Type, data any) error {
	for _, h := range e.hooks[t] {
		run := func() error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[PANIC] hook=%s err=%v\n", h.Name, r)
				}
			}()
			return h.Fn(ctx, data)
		}

		if h.Mode == Async {
			go run()
			continue
		}

		if err := run(); err != nil && h.MustSucceed {
			return err
		}
	}
	return nil
}
