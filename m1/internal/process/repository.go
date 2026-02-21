package process

import (
	"context"
	"database/sql"
	"fmt"
)

type RepositoryInterface interface {
	Register(ctx context.Context, dto RegisterDTO) error
	UpdatePublishedStatus(ctx context.Context, traceID string) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Register(ctx context.Context, dto RegisterDTO) error {
	query := RegisterQuery()
	_, err := r.db.ExecContext(ctx, query,
		dto.TraceID,
		dto.Payload,
		dto.ByteSize,
		dto.TotalCharacters,
	)
	if err != nil {
		return fmt.Errorf("failed to persist record: %w", err)
	}

	return nil
}

func (r *Repository) UpdatePublishedStatus(ctx context.Context, traceID string) error {
	query := UpdatePublishedQuery()
	_, err := r.db.ExecContext(ctx, query,
		traceID,
	)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	return nil
}
