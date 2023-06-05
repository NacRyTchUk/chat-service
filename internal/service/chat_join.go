package service

import (
	"chat-service/internal/dto"
	"time"
)

func (c chatService) ChatJoin(name string, chatId int64) (dto.Chat, []dto.Message, error) {
	user, err := c.dao.NewUserQuery().GetId(name)
	if err != nil {
		return dto.Chat{}, nil, err
	}

	_, err = c.dao.NewUserChatsQuery().Create(user.Id, chatId)
	if err != nil {
		return dto.Chat{}, nil, err
	}
	chat, err := c.dao.NewChatQuery().GetName(chatId)
	if err != nil {
		return dto.Chat{}, nil, err
	}
	msgs, err := c.dao.NewMessageQuery().GetChatMessages(chatId, time.Minute)
	if err != nil {
		return dto.Chat{}, nil, err
	}

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

	c.ChatHandler(dto.FormerMessage{
		ChatId:     chatId,
		SenderName: name,
		Text:       "Has joined the chat",
	})

	return dto.Chat(*chat), messages, nil
}
