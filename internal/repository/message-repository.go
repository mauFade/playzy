package repository

import (
	"database/sql"

	"github.com/mauFade/playzy/internal/model"
)

type MessageRepositoryInterface interface {
	Create(m model.Message) error
	List(fstUserId, scdUserId string, limit, offset int) ([]model.Message, error)
	SetMessagesIsRead(userId, otherUserID string) error
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

		CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);
		CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages(receiver_id);
		CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
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

func (r *MessageRepository) List(fstUserId, scdUserId string, limit, offset int) ([]model.Message, error) {
	query := `
		SELECT id, content, user_id, receiver_id, created_at, is_read
		FROM messages
		WHERE (user_id = $1 AND receiver_id = $2) OR (user_id = $2 AND receiver_id = $1)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, fstUserId, scdUserId, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := []model.Message{}

	for rows.Next() {
		var msg model.Message

		err := rows.Scan(&msg.ID,
			&msg.Content,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Timestamp,
			&msg.IsRead)

		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *MessageRepository) SetMessagesIsRead(userId, otherUserID string) error {
	_, err := r.db.Exec(`
			UPDATE messages
			SET is_read = true
			WHERE receiver_id = $1 AND user_id = $2 AND is_read = false
		`, userId, otherUserID)

	if err != nil {
		return err
	}

	return nil
}
