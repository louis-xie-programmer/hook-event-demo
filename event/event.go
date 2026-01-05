package event

type OrderCreatedEvent struct {
	OrderID string
	UserID  string
	Amount  int64
}
