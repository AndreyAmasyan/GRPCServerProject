package broker

import "test/internal/data"

// MessageBroker ...
type MessageBroker interface {
	Produce(user data.User) error
	Close() error
}
