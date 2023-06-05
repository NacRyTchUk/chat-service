package service

import (
	"chat-service/internal/dto"
	"log"
	"time"
)

func (c chatService) ChatJoin(name string, chatId int64) (dto.Chat, []dto.Message, error) {
	log.Println(1, name, chatId)
	user, err := c.dao.NewUserQuery().GetId(name)
	if err != nil {
		return dto.Chat{}, nil, err
	}
	log.Println(2, user)

	_, err = c.dao.NewUserChatsQuery().Create(user.Id, chatId)
	if err != nil {
		return dto.Chat{}, nil, err
	}
	log.Println(3, err)
	chat, err := c.dao.NewChatQuery().GetName(chatId)
	if err != nil {
		return dto.Chat{}, nil, err
	}
	log.Println(4, chat)
	msgs, err := c.dao.NewMessageQuery().GetChatMessages(chatId, time.Minute)
	if err != nil {
		return dto.Chat{}, nil, err
	}
	log.Println(5, msgs)

	var messages []dto.Message
	for _, msg := range msgs {
		name, _ := c.dao.NewUserQuery().GetName(msg.SenderId)
		messages = append(messages, dto.Message{
			Id:         msg.Id,
			ChatId:     msg.ChatId,
			SenderName: name.Name,
			Text:       msg.Text,
			Timestamp:  msg.CreatedAt.Unix(),
		})
	}

	return dto.Chat(*chat), messages, nil
}
