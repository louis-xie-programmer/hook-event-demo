package event

import "sync"

type Handler func(e any)

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewBus() *Bus {
	return &Bus{handlers: make(map[string][]Handler)}
}

func (b *Bus) Subscribe(topic string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], h)
}

func (b *Bus) Publish(topic string, e any) {
	b.mu.RLock()
	hs := b.handlers[topic]
	b.mu.RUnlock()

	for _, h := range hs {
		go h(e)
	}
}
