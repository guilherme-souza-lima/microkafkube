package process

import (
	"context"
	"database/sql"
	"fmt"
)

type OutboxRepository interface {
	GetPendingRegistrations(ctx context.Context) ([]RegistrationDTO, error)
	UpdatePublishedStatus(ctx context.Context, traceID string) error
}

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

func (r *Repository) GetPendingRegistrations(ctx context.Context) ([]RegistrationDTO, error) {
	query := GetByPublishedFalseQuery()
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending registrations: %w", err)
	}
	defer rows.Close()

	var registrations []RegistrationDTO
	for rows.Next() {
		var reg RegistrationDTO
		err := rows.Scan(
			&reg.TraceID,
			&reg.Payload,
			&reg.ByteSize,
			&reg.TotalCharacters,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan registration: %w", err)
		}
		registrations = append(registrations)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}
