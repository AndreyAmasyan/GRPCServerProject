package storage

import (
	"context"
	"test/internal/data"
)

// Storage ...
type Storage interface {
	AddUser(ctx context.Context, user data.User) error
	FetchAllUsers(ctx context.Context) ([]data.User, error)
	DeleteUser(ctx context.Context, id int) error
	Init(ctx context.Context) error
	Close() error
}

// Cache ...
type Cache interface {
	SetAllUsers(ctx context.Context, list []data.User) error
	IsExists(ctx context.Context) (bool, error)
	GetAllUsers(ctx context.Context) ([]data.User, error)
	Close() error
}

// Log ...
type Log interface {
	Init(ctx context.Context) error
	Close() error
}
