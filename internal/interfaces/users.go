package interfaces

import "context"

type User interface {
	CreateUser(ctx context.Context)
	GetUser(ctx context.Context)
	UpdateUser(ctx context.Context)
	DeleteUser(ctx context.Context)
}
