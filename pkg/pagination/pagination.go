package pagination

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/rotisserie/eris"
	"github.com/voxtmault/mentoring/library-project/pkg/utility"
)

type PaginationMetadata struct {
	TotalRecords uint `json:"total_records"`
	TotalPages   uint `json:"total_pages"`
	CurrentLimit uint `json:"page_size"`
	CurrentPage  uint `json:"current_page"`
}

type PaginationFilter struct {
	Limit      uint `query:"page_size" validate:"required,number,min=1"`
	PageNumber uint `query:"page_number" validate:"required,number,min=1"`
}

func ProcessPaginationRequest(ctx context.Context, con *sql.DB, statement string, args []interface{}, metadata *PaginationMetadata) error {
	// Update the pagination metadata
	args = utility.RemoveLastTwoItems(args)
	if err := con.QueryRowContext(ctx, statement, args...).Scan(&metadata.TotalRecords); err != nil {
		slog.Error("failed to get queue data count", "error", err)
		return eris.Wrap(err, "failed to get queue data count")
	}

	metadata.TotalPages = (metadata.TotalRecords + metadata.CurrentLimit - 1) / metadata.CurrentLimit

	return nil
}
