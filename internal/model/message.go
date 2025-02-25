package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID `json:"id"`
	Content    string    `json:"content"`
	SenderID   string    `json:"senderId"`
	ReceiverID string    `json:"receiverId"`
	Timestamp  time.Time `json:"timestamp"`
	IsRead     bool      `json:"isRead"`
}

func NewMessage(id uuid.UUID,
	content,
	senderID,
	receiverID string,
	timestamp time.Time,
	isRead bool) *Message {
	return &Message{
		ID:         id,
		Content:    content,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Timestamp:  timestamp,
		IsRead:     isRead,
	}
}

func (m *Message) GetID() uuid.UUID {
	return m.ID
}

func (m *Message) GetContent() string {
	return m.Content
}

func (m *Message) SetContent(content string) {
	m.Content = content
}

func (m *Message) GetSenderID() string {
	return m.SenderID
}

func (m *Message) SetSenderID(senderID string) {
	m.SenderID = senderID
}

func (m *Message) GetReceiverID() string {
	return m.ReceiverID
}

func (m *Message) SetReceiverID(receiverID string) {
	m.ReceiverID = receiverID
}

func (m *Message) GetTimestamp() time.Time {
	return m.Timestamp
}

func (m *Message) SetTimestamp(timestamp time.Time) {
	m.Timestamp = timestamp
}

func (m *Message) GetIsRead() bool {
	return m.IsRead
}

func (m *Message) SetIsRead(isRead bool) {
	m.IsRead = isRead
}
