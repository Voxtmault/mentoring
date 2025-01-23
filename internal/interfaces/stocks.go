package interfaces

import "context"

type Stock interface {
	AddStock(ctx context.Context)
	GetStock(ctx context.Context)
	UpdateStock(ctx context.Context)
	DeleteStock(ctx context.Context)
}
