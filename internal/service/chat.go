package service

import (
	"chat-service/internal/broker"
	"chat-service/internal/dto"
	"chat-service/internal/repository"
)

type ChatService interface {
	ChatHandler(message dto.FormerMessage) error
	ChatList() ([]dto.Chat, error)
	ChatJoin(name string, chatId int64) (dto.Chat, []dto.Message, error)
	NewUser(name string) (dto.User, error)
	GetUsers(chatId int64) ([]dto.User, error)
}

type chatService struct {
	dao    repository.DAO
	broker broker.MessageBroker
}

func NewChatService(dao repository.DAO,
	broker broker.MessageBroker,
) ChatService {
	return &chatService{
		dao:    dao,
		broker: broker,
	}
}
