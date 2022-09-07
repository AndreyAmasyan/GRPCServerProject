package service

import (
	"context"
	"test/internal/dto"
)

// Processor ...
type Processor interface {
	AddUser(ctx context.Context, user dto.User) error
	FetchAllUsers(ctx context.Context) ([]dto.User, error)
	DeleteUser(ctx context.Context, id int) error
}
