package interfaces

import "context"

type Transaction interface {
	CreateTransaction(ctx context.Context)
	GetTransaction(ctx context.Context)
	UpdateTransaction(ctx context.Context)
	DeleteTransaction(ctx context.Context)
}
