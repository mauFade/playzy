package message

import (
	"log"
	"time"

	"github.com/mauFade/playzy/internal/repository"
)

type ListUsersMessagesUseCase struct {
	mr repository.MessageRepositoryInterface
}

type ListUsersMessagesRequest struct {
	UserID      string
	OtherUserID string
	Limit       int
	Offset      int
}

type MessageResponse struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	SenderID   string    `json:"senderId"`
	ReceiverID string    `json:"receiverId"`
	Timestamp  time.Time `json:"timestamp"`
	IsRead     bool      `json:"isRead"`
	IsMine     bool      `json:"isMine"`
}

func NewListUsersMessagesUseCase(mr repository.MessageRepositoryInterface) *ListUsersMessagesUseCase {
	return &ListUsersMessagesUseCase{
		mr: mr,
	}
}

func (uc *ListUsersMessagesUseCase) Execute(data *ListUsersMessagesRequest) ([]MessageResponse, error) {
	ms, err := uc.mr.List(data.UserID, data.OtherUserID, data.Limit, data.Offset)

	if err != nil {
		return nil, err
	}

	var messages []MessageResponse
	for _, m := range ms {
		var msg MessageResponse

		msg.ID = m.ID.String()
		msg.Content = m.Content
		msg.SenderID = m.SenderID
		msg.ReceiverID = m.ReceiverID
		msg.Timestamp = m.Timestamp
		msg.IsRead = m.IsRead
		msg.IsMine = (m.SenderID == data.UserID)

		messages = append(messages, msg)
	}

	// Atualizar mensagens como lidas (se o usuário for o destinatário)
	// go func() {
	err = uc.mr.SetMessagesIsRead(data.UserID, data.OtherUserID)

	log.Printf("Erro ao marcar mensagens como lidas: %v", err)
	if err != nil {
		log.Printf("Erro ao marcar mensagens como lidas: %v", err)
	}
	// }()

	return messages, nil
}
