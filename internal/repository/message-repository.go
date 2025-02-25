package repository

import (
	"database/sql"

	"github.com/mauFade/playzy/internal/model"
)

type MessageRepositoryInterface interface {
	Create(m model.Message) error
}

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(d *sql.DB) *MessageRepository {
	r := &MessageRepository{
		db: d,
	}

	r.db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				content TEXT NOT NULL,
				user_id UUID NOT NULL REFERENCES users(id),
				receiver_id UUID NOT NULL REFERENCES users(id),
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				is_read BOOLEAN DEFAULT false
		);

		CREATE INDEX idx_messages_user_id ON messages(user_id);
		CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);
		CREATE INDEX idx_messages_created_at ON messages(created_at);
	`)

	return r
}

func (r *MessageRepository) Create(m model.Message) error {
	var messageID string

	err := r.db.QueryRow(`
        INSERT INTO messages (content, user_id, receiver_id, created_at, is_read)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, m.Content, m.SenderID, m.ReceiverID, m.Timestamp, m.IsRead).Scan(&messageID)

	if err != nil {
		return err
	}

	return nil
}
