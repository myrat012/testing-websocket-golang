package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewInMemoryMessageRepository(d *pgxpool.Pool, cx context.Context) MessageRepository {
	return MessageRepository{
		db:  d,
		ctx: cx,
	}
}

func (m *MessageRepository) SaveMessage(message string) error {
	_, err := m.db.Exec(m.ctx, "INSERT INTO messages (message) VALUES ($1)", message)
	return err
}
