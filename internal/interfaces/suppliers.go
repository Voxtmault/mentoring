package interfaces

import "context"

type Supplier interface {
	CreateSupplier(ctx context.Context)
	GetSupplier(ctx context.Context)
	UpdateSupplier(ctx context.Context)
	DeleteSupplier(ctx context.Context)
}
